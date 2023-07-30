package digital

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"image"
	"image/color"
	"testing"
)

func TestSeconds(t *testing.T) {
	d := display.NewTestRGBA(64, 20, t.Name())
	c := SecondClock(text.Large, color.White)
	c(d.Image(), image.Point{})
	d.Update()
}
