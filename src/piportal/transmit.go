package piportal

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net"
	"os"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/thiswillgowell/light-controller/ratetracker"

	"golang.org/x/time/rate"
)

type Connection struct {
	conn           net.Conn
	address        string
	counter        *int64
	packetLimitter *rate.Limiter
	updateLimitter *rate.Limiter
	tracker        ratetracker.Tracker
	reconnect      func()
}

const (
	chunkSize = 4096
)

func (c *Connection) WriteFrame(img image.Image) {

	if err := c.updateLimitter.Wait(context.Background()); err != nil {
		zap.S().Errorw("failed to wait on update limiter", zap.Error(err))
		return
	}

	pixelData := getPixelData(img)
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	//// Prepare the length of the pixel data
	//dataLen := uint32(len(pixelData))
	//var err error
	//// Send the length of the pixel data
	//err = binary.Write(c.conn, binary.BigEndian, dataLen)
	//if err != nil {
	//	zap.S().Errorw("client is disconnected")
	//	return
	//}

	// Send the pixel data in chunks
	totalSent := 0
	for totalSent < len(pixelData) {
		end := totalSent + chunkSize
		if end > len(pixelData) {
			end = len(pixelData)
		}
		chunk := pixelData[totalSent:end]

		n, err := c.conn.Write(chunk)
		if err != nil {
			zap.S().Errorw("Connection closed. Waiting for a new connection...")
			return
		}

		totalSent += n
	}
	atomic.AddInt64(c.counter, 1)
}

func getPixelData(img image.Image) []byte {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	pixelData := make([]byte, width*height*3)

	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			pixelData[index] = c.R
			pixelData[index+1] = c.G
			pixelData[index+2] = c.B
			index += 3
		}
	}

	return pixelData
}

func NewUDPClient(address string) (*Connection, error) {

	c := &Connection{
		counter:        new(int64),
		packetLimitter: rate.NewLimiter(rate.Every(time.Nanosecond*100), 1),
		updateLimitter: rate.NewLimiter(150, 1),
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

func NewTCPClient(address string) (*Connection, error) {

	c := &Connection{
		counter: new(int64),
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
