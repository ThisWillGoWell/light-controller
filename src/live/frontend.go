package live

import (
	_ "embed"
	"fmt"
	"github.com/thiswillgowell/light-controller/src/display"
	"go.uber.org/zap"
	"log"
	"net/http"
)

//go:embed liveView.html
var liveViewHTML []byte

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(liveViewHTML); err != nil {
		zap.L().Error("failed to write html", zap.Error(err))
	}
}

func RunServer(subDisplay *display.SubscriptionDisplay) {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(subDisplay, w, r)
	})
	fmt.Println("running")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
