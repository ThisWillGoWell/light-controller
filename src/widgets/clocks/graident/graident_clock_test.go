package graident

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"image"
	"testing"
)

func TestClock(t *testing.T) {
	size := 64
	d := display.NewTestRGBA(size, size, t.Name())
	updateFn := Clock(Config{
		Size:   size,
		Colors: CoolBluesColorSet(),
	})
	updateFn(d.Image(), image.Point{})
	d.Update()
}
