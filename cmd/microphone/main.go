package main

import (
	"github.com/thiswillgowell/light-controller/src/audio/inputs/microphone"
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/live"
	"golang.org/x/image/colornames"
	"image/color"
)

func main() {
	device := microphone.Microphone

	device.Start()

	d := display.NewRGBA(64, 64)
	go live.RunServer(display.NewSubscription(d))
	music.NewVuMeter(music.Params{
		BarWidth: 1,
		Display:  d,
		BarColors: func(i int, i2 int) []color.Color {
			colors := make([]color.Color, i2)
			for i := range colors {
				colors[i] = colornames.Orange
			}
			return colors
		},
		Channel: specturm.LeftChannel,
	})
	select {}
}
