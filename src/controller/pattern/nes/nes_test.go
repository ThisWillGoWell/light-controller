package nes

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/live"
	"github.com/thiswillgowell/light-controller/src/piportal"
	"testing"
	"time"
)
import "github.com/fogleman/nes/nes"
import "github.com/stretchr/testify/require"

func TestNES(t *testing.T) {
	console, err := nes.NewConsole("/Users/will/workspace/light-controller/src/controller/pattern/nes/Donkey Kong (World) (Rev A).nes")
	require.NoError(t, err)

	sub := display.NewSubscription(piportal.Fireplace)

	go live.RunServer(sub)

	tick := time.NewTicker(time.Second / 60)
	for {
		<-tick.C
		console.StepFrame()
		drawFrame(sub, console)
	}
}
