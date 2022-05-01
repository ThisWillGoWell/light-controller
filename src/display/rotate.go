package display

import (
	"image"
	"image/color"
	"image/draw"
)

type Rotation int

const (
	NoRotation Rotation = iota
	Clockwise
	CounterClockwise
	OneEighty
	MirrorAcrossY
)

type TransposeFunc func(int, int) (int, int)

func noRotation(x, y int) (int, int) {
	return x, y
}

func rotateCounterClockwise(x, y int) (int, int) {
	return y, x
}

func rotate180(inX, inY int) func(int, int) (int, int) {
	return func(x int, y int) (int, int) {
		return inX - x - 1, inY - y - 1
	}
}

//
func mirrorAcrossY(maxY int) TransposeFunc {
	return func(x int, y int) (int, int) {
		return x, maxY - 1 - y
	}
}

type VirtualDisplayRotation struct {
	X                 int
	Y                 int
	underlyingDisplay Display
	transposeFunc     TransposeFunc
	boundingBox       image.Rectangle
}

func (v *VirtualDisplayRotation) Image() draw.Image {
	return v
}

func (v *VirtualDisplayRotation) UpdateImage(src image.Image) {
	draw.Draw(v, v.Bounds(), src, image.Point{}, draw.Src)
	v.Update()
}

func (v *VirtualDisplayRotation) Update() {
	v.underlyingDisplay.Update()
}

func (v *VirtualDisplayRotation) ColorModel() color.Model {
	return color.RGBAModel
}

func (v *VirtualDisplayRotation) Bounds() image.Rectangle {
	return v.boundingBox
}

func (v *VirtualDisplayRotation) At(x, y int) color.Color {
	x, y = v.transposeFunc(x, y)
	return v.underlyingDisplay.Image().At(x, y)
}

func (v *VirtualDisplayRotation) Set(x, y int, c color.Color) {
	x, y = v.transposeFunc(x, y)
	v.underlyingDisplay.Image().Set(x, y, c)
}

func NewRotation(d Display, rotationType Rotation) *VirtualDisplayRotation {

	dRotate := &VirtualDisplayRotation{
		underlyingDisplay: d,
	}
	switch rotationType {
	case NoRotation:
		dRotate.transposeFunc = noRotation
		dRotate.Y = d.Image().Bounds().Size().Y
		dRotate.X = d.Image().Bounds().Size().X
	case Clockwise:
		dRotate.Y = d.Image().Bounds().Size().X
		dRotate.X = d.Image().Bounds().Size().Y
		dRotate.transposeFunc = rotateCounterClockwise
	case OneEighty:
		dRotate.X = d.Image().Bounds().Size().X
		dRotate.Y = d.Image().Bounds().Size().Y
		dRotate.transposeFunc = rotate180(d.Image().Bounds().Size().X, d.Image().Bounds().Size().Y)
	case CounterClockwise:
		dRotate.Y = d.Image().Bounds().Size().X
		dRotate.X = d.Image().Bounds().Size().Y
		dRotate.transposeFunc = rotateCounterClockwise
	case MirrorAcrossY:
		dRotate.Y = d.Image().Bounds().Size().Y
		dRotate.X = d.Image().Bounds().Size().X
		dRotate.transposeFunc = mirrorAcrossY(dRotate.Y)
	}
	dRotate.boundingBox = image.Rect(0, 0, dRotate.X, dRotate.Y)
	return dRotate
}
