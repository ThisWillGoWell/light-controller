package live

import (
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

func runConnection(v *display.SubscriptionDisplay, conn *websocket.Conn) {
	connID := conn.RemoteAddr().String() + time.Now().String()

	imageChannel := make(chan image.Image)
	v.RegisterSubscription(connID, imageChannel)
	for img := range imageChannel {
		rgba, ok := img.(*image.RGBA)
		if !ok {
			panic("Hey i dont know what do to for a non rgba image")
		}
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
		_ = rgba
		if err := conn.WriteMessage(websocket.TextMessage, []byte("hello-world")); err != nil {
			zap.L().Warn("closing connection because of error", zap.Error(err))
			fmt.Println(err.Error())
			close(imageChannel)
		}
	}
	v.RemoveSubscription(connID)
	if err := conn.Close(); err != nil {
		zap.L().Error("failed to close connection")
	}
}

func serveWs(subscription *display.SubscriptionDisplay, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go runConnection(subscription, conn)
}
