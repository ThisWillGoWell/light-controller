package display

import (
	"image"

	"golang.org/x/image/draw"
)

type Display interface {
	Image() draw.Image
	UpdateImage(image image.Image)
	Update()
}
