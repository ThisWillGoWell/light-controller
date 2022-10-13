package fire

import (
	"github.com/thiswillgowell/light-controller/src/piportal"
	"testing"
	"time"
)

func TestFire(t *testing.T) {
	//pattern := NewFlamePattern(display.NewTestRGBA(60, 30, "fire"), 75)
	pattern := NewFlamePattern(piportal.Fireplace, 75)
	tick := time.NewTicker(time.Second / 75)
	for {
		<-tick.C
		pattern.step()
	}
}
