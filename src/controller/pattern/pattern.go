package pattern

import (
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"image"
)

type Pattern interface {
	ImageChannel() chan *image.Image

	Start()
	Stop()
}

type MusicPattern interface {
	ProcessSpectrum(specturm.FrequencySpectrum) error
}
