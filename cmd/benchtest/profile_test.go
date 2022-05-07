package benchtest

import (
	"testing"
	"time"

	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/piportal"
)

func TestProfile(t *testing.T) {
	daisyDevice, err := daisy.Init()
	if err != nil {
		panic(err)
	}

	p1, err := piportal.NewMatrix("192.168.1.53:8080", piportal.TopLeft)
	if err != nil {
		panic(err)
	}
	//p2, err := piportal.NewMatrix("192.168.1.83:8080", piportal.TopRight)
	//if err != nil {
	//	panic(err)
	//}

	//p := display.NewRotation(display.NewMultiDisplay(display.ArrangementVertical, display.NewRotation(p2, display.MirrorAcrossY), p1), display.MirrorAcrossY)

	go func() {
		music.CenterHollowVUBarDouble(daisyDevice, p1, 1)
	}()
	<-time.After(time.Second * 3)
}
