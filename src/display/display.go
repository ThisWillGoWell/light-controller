package display

import "github.com/thiswillgowell/light-controller/color"

type Display interface {
	Rows() int
	Cols() int
	Clear()
	GetPixel(row, col int) color.Color
	SetPixel(row, col int, c color.Color)
	Send() error
}

func