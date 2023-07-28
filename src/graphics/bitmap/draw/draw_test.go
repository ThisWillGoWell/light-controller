package draw

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"golang.org/x/image/colornames"
	"testing"
)

func TestT(*testing.T) {
	img := display.NewTestRGBA(100, 100, "test")
	targetImg := img.Image()
	drawLine(targetImg, colornames.Purple, 4, 0, 0, 100, 0)
	drawLine(targetImg, colornames.Red, 4, 100, 0, 100, 100)
	drawLine(targetImg, colornames.Gold, 4, 100, 100, 0, 100)
	drawLine(targetImg, colornames.Blue, 4, 0, 100, 0, 0)

	drawLine(targetImg, colornames.Black, 4, 0, 0, 100, 100)
	img.Update()

}
