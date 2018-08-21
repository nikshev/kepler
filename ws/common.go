package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lastexile/kepler"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3 * 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// default connection retry interval
	connectionRetryInterval = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnectionFactoryFunc returns web socket open connection function
type ConnectionFactoryFunc func() (*websocket.Conn, error)

// SendTextMessage data over opened conn
func SendTextMessage(conn *websocket.Conn, payload []byte) error {
	return write(conn, websocket.TextMessage, payload)
}

// write writes a message with the given message type and payload.
func write(conn *websocket.Conn, mt int, payload []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(mt, payload)
}

// DialConnection returns based on Default dialer connectionFactory
func DialConnection(path string) func() (*websocket.Conn, error) {
	return func() (*websocket.Conn, error) {
		conn, _, err := websocket.DefaultDialer.Dial(path, nil)
		return conn, err
	}
}

// ServeConnection returns based on Default web socket serve connectionFactory
func ServeConnection(w http.ResponseWriter, r *http.Request) func() (*websocket.Conn, error) {
	return func() (*websocket.Conn, error) {
		return upgrader.Upgrade(w, r, nil)
	}
}

// JsonValue return default json marshal data
func JsonValue(msg kepler.Message) ([]byte, error) {
	return json.Marshal(msg.Value())
}