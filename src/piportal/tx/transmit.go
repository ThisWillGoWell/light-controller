package tx

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/thiswillgowell/light-controller/color"
	"golang.org/x/time/rate"
)

type Connection struct {
	conn           net.Conn
	counter        *int64
	packetLimitter *rate.Limiter
	updateLimitter *rate.Limiter
}

func (c *Connection) WriteFrame(frame [][]color.Color) error {
	if err := c.updateLimitter.Wait(context.Background()); err != nil {
		return err
	}
	packet := make([]byte, len(frame[0])*3+1)
	for r, row := range frame {
		packet[0] = byte(uint(r))
		for c, colorObj := range row {
			packet[c*3+1] = colorObj.R
			packet[c*3+2] = colorObj.G
			packet[c*3+3] = colorObj.B
		}
		if _, err := c.conn.Write(packet); err != nil {
			return err
		}
	}
	return nil
}

func NewClient(address string) (*Connection, error) {

	c := &Connection{
		counter:        new(int64),
		packetLimitter: rate.NewLimiter(rate.Every(time.Nanosecond*100), 1),
		updateLimitter: rate.NewLimiter(100, 1),
	}
	var err error
	if c.conn, err = net.Dial("udp", address); err != nil {
		return nil, err
	}

	go func() {
		t := time.NewTicker(time.Second)
		for {
			<-t.C
			fmt.Printf("%d requests/second\n", *c.counter)
			atomic.SwapInt64(c.counter, 0)
		}

	}()

	return c, nil
}
