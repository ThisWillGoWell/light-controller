package display

import (
	"github.com/thiswillgowell/light-controller/color"
)

type Rotation int

const (
	NoRotation Rotation = iota
	Clockwise
	CounterClockwise
	OneEighty
)

type TransposeFunc func(int, int) (int, int)

func noRotation(inR, inC int) (int, int) {
	return inR, inC
}

func rotateCounterClockwise(inR, inC int) (int, int) {
	return inC, inR
}

func rotate180(inRows, inCols int) func(int, int) (int, int) {
	return func(inR int, inC int) (int, int) {
		return inRows - inR - 1, inCols - inC - 1
	}
}

type VirtualDisplayRotation struct {
	numRows           int
	numCols           int
	underlyingDisplay Display
	transposeFunc     TransposeFunc
}

func NewRotation(d Display, rotationType Rotation) Display {

	dRotate := &VirtualDisplayRotation{
		underlyingDisplay: d,
	}
	switch rotationType {
	case NoRotation:
		dRotate.transposeFunc = noRotation
		dRotate.numRows = d.Rows()
		dRotate.numCols = d.Cols()
	case Clockwise:
		panic("not impl")
	case OneEighty:
		dRotate.numRows = d.Rows()
		dRotate.numCols = d.Cols()
		dRotate.transposeFunc = rotate180(d.Rows(), d.Cols())
	case CounterClockwise:
		dRotate.numRows = d.Cols()
		dRotate.numCols = d.Rows()
		dRotate.transposeFunc = rotateCounterClockwise
	}
	return dRotate
}

func (v VirtualDisplayRotation) Rows() int {
	return v.numRows
}

func (v VirtualDisplayRotation) Cols() int {
	return v.numCols
}

func (v VirtualDisplayRotation) Clear() {
	v.underlyingDisplay.Clear()
}

func (v VirtualDisplayRotation) GetPixel(row, col int) color.Color {
	r, c := v.transposeFunc(row, col)
	return v.underlyingDisplay.GetPixel(r, c)
}

func (v VirtualDisplayRotation) SetPixel(row, col int, c color.Color) {
	r, col := v.transposeFunc(row, col)
	v.underlyingDisplay.SetPixel(r, col, c)
}

func (v VirtualDisplayRotation) Send() error {
	return v.underlyingDisplay.Send()
}
