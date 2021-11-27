package pattern

import "github.com/thiswillgowell/light-controller/color_2"

type Pattern interface {
	GetNextValue() [][]color_2.Color
}

type MusicPattern struct {
}
