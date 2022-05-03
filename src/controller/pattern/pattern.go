package pattern

import (
	"image"
)

type Pattern interface {
	ImageChannel() chan *image.Image
	Start()
	Stop()
}

type MusicPattern struct {
}
