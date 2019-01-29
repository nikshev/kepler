package rmq

import (
	"time"

	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// NewSink - return generic amqp publisher sink
func NewSink(connFactory ConnectionFactoryFunc, queue QueueOptions, formatter kepler.MarshallerFunc) (kepler.Sink, error) {

	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error
	var data []byte

	return kepler.NewSink(func(msg kepler.Message) {
		for {
			if conn == nil || ch == nil {
				log.Infof("Connection or channel is closed %v - %v\n", conn, ch)
				for {

					conn, ch, err = initExchangeChannelQueue(connFactory, queue)

					if err != nil {
						log.Errorf("Failed to init channel: %v\n", err)
						conn, ch = resetConnection(conn, ch)

						time.Sleep(connectionRetryInterval)
						continue
					}

					if len(queue.Name) != 0 {
						err = ch.QueueBind(
							queue.Name,         // name of the queue
							queue.RoutingKey,   // bindingKey
							queue.ExchangeName, // sourceExchange
							queue.NoWait,       // noWait
							nil,                // arguments
						)
					}

					if err != nil {
						log.Errorf("Failed to binding a queue: %v\n", err)
						conn, ch = resetConnection(conn, ch)
						continue
					}
					break
				}
			}

			data, err = formatter(msg)
			err = ch.Publish(
				queue.ExchangeName, // exchange
				queue.RoutingKey,   // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
				})

			if err != nil {
				log.Errorf("Failed to publish: %v\n", err)
				conn, ch = resetConnection(conn, ch)
				continue
			}
			log.Printf("Published: %v - %v", string(data), err)
			break
		}
	}), nil
}
