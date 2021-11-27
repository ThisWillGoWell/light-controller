package benchtest

import (
	"testing"
	"time"

	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
)

func TestProfile(t *testing.T) {
	daisyDevice, err := daisy.InitDaisy()
	if err != nil {
		panic(err)
	}

	matrixRightDisplay, err := piportal.NewMatrix("192.168.1.53:8080")
	if err != nil {
		panic(err)
	}

	matrixLeftDisplay, err := piportal.NewMatrix("192.168.1.63:8080")
	if err != nil {
		panic(err)
	}

	combindedDispaly := display.NewMultiDisplay(display.ArrangementHorizontal,
		display.NewRotation(matrixLeftDisplay, display.OneEighty),
		matrixRightDisplay,
	)

	go func() {
		music.CenterHollowVUBarDouble(daisyDevice, combindedDispaly, 1)
	}()
	<-time.After(time.Second * 3)
}
