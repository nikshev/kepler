package kepler

import (
	"context"
	"fmt"
	"testing"
)

func TestSpringFanout(t *testing.T) {
	s := NewSpring(func(ctx context.Context, c chan<- Message) {
		for i := 0; i < 10; i++ {
			c <- NewMessage("range", i)
		}
	})

	t1 := NewSink(func(m Message) {
		fmt.Println("t1: " + m.String())
	})

	t2 := NewSink(func(m Message) {
		fmt.Println("t2: " + m.String())
	})

	s.LinkTo(".", t2, Allways)
	s.LinkTo(".", t1, Allways)
}
