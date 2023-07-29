package weather

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"image"
	"testing"
)

func TestVerticalFullInfoDrawWeather(t *testing.T) {
	d := display.NewTestRGBA(64, 96, t.Name())
	VerticalFullInfoDraw(d.Image(), image.Point{})
	d.Update()
}
