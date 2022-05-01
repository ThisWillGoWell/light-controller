package piportal

import (
	"context"
	"fmt"
	"image"
	"net"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/thiswillgowell/light-controller/ratetracker"

	"golang.org/x/time/rate"
)

type Connection struct {
	conn           net.Conn
	mode           PortalMode
	address        string
	counter        *int64
	packetLimitter *rate.Limiter
	updateLimitter *rate.Limiter
	tracker        ratetracker.Tracker
	reconnect      func()
}

func (c *Connection) WriteFrame(image *image.RGBA) {
	if err := c.updateLimitter.Wait(context.Background()); err != nil {
		zap.S().Errorw("failed to wait on update limiter", zap.Error(err))
		return
	}
	if c.conn == nil {
		return
	}
	atomic.AddInt64(c.counter, 1)

	maxX := image.Rect.Max.X
	maxY := image.Rect.Max.Y

	packet := make([]byte, maxX*3+1)

	for y := 0; y < maxY; y++ {
		imageY := y
		if c.mode == Right {
			if y < 32 {
				imageY = 31 - y
			} else if y >= 64 {
				imageY = 64 + (95 - y)
			}
		}
		if y >= 32 && y < 64 {
			imageY = 32 + 63 - y
		}
		packet[0] = byte(y)
		for x := 0; x < maxX; x++ {
			imageX := x
			if c.mode == Left && y >= 32 && y < 64 {
				imageX = 63 - x
			}
			r, g, b, _ := image.RGBAAt(imageX, imageY).RGBA()
			packet[x*3+1] = uint8(r)
			packet[x*3+2] = uint8(g)
			packet[x*3+3] = uint8(b)
		}
		if _, err := c.conn.Write(packet); err != nil {
			zap.S().Warnw("error when writing to display, starting the reconnect process...", zap.Error(err))
			_ = c.conn.Close()
			c.conn = nil
			go c.reconnect()
		}
	}
}

func NewUDPClient(address string) (*Connection, error) {

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

func NewTCPClient(address string, mode PortalMode) (*Connection, error) {

	c := &Connection{
		counter: new(int64),
		mode:    mode,
		address: address,
		//packetLimitter: rate.NewLimiter(rate.Every(time.Nanosecond*100), 1),
		updateLimitter: rate.NewLimiter(100, 1),
	}

	var err error
	if c.conn, err = net.Dial("tcp", address); err != nil {
		zap.S().Warnw("could not connect to client, starting reconnect", zap.Error(err))

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

func (c *Connection) reconnectTCP(address string) func() {
	return func() {
		var err error
		if c.conn, err = net.Dial("tcp", address); err == nil {
			zap.S().Infow("reconnected to portal", "address", address)
			return
		}
		<-time.After(time.Second * 15)
	}
}
