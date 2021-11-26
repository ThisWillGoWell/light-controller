package canvas

import (
	"image"

	"github.com/thiswillgowell/light-controller/src/display"
)

type BaseCanvas struct {
	display.Display
	Image image.RGBA
}

type Canvas struct {
}

func NewBaseCanvas() {
	return
}
