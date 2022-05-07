package live

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/thiswillgowell/light-controller/src/display"
	"go.uber.org/zap"
	"image"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
}

type message struct {
	Payload string `json:"payload"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

func encodeImage(image image.Image) []byte {
	max := image.Bounds().Max
	payload := make([]byte, (max.X+1)*(max.Y+1)*3)
	for x := 0; x < max.X; x++ {
		for y := 0; y < max.Y; y++ {
			r, g, b, _ := image.At(x, y).RGBA()
			payload[(x+y*max.X)*3] = uint8(r)
			payload[(x+y*max.X)*3+1] = uint8(g)
			payload[(x+y*max.X)*3+2] = uint8(b)
		}
	}
	val, _ := json.Marshal(message{
		Payload: hex.EncodeToString(payload),
		Width:   max.X,
		Height:  max.Y,
	})
	return val
}

func runConnection(v *display.SubscriptionDisplay, conn *websocket.Conn) {
	connID := conn.RemoteAddr().String() + time.Now().String()
	imageChannel := v.RegisterSubscription(connID)
	defer func() {
		if err := conn.Close(); err != nil {
			zap.L().Error("failed to close connection", zap.Error(err))
		}
	}()
	defer v.RemoveSubscription(connID)
	for range time.Tick(time.Second / 30) {
		img := <-imageChannel
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, encodeImage(img)); err != nil {
			zap.L().Warn("closing connection because of error", zap.Error(err))
			fmt.Println(err.Error())
			return
		}
	}
}

func serveWs(subscription *display.SubscriptionDisplay, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go runConnection(subscription, conn)
}
