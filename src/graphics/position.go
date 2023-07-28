package graphics

import (
	"image"
	"image/draw"
)

func DawImageAtMidpointX(yPos int, source image.Image, dest draw.Image) {
	targetPos := image.Point{
		X: int(float64(dest.Bounds().Size().X)/2.0 - float64(source.Bounds().Size().X)/2.0),
		Y: yPos,
	}
	endPos := targetPos.Add(source.Bounds().Size())
	draw.Draw(dest, image.Rect(targetPos.X, targetPos.Y, endPos.X, endPos.Y), source, image.Point{}, draw.Over)
}

func PositionImageAround(pos image.Point, source image.Image, dest draw.Image) {
	endPos := pos.Add(pos.Div(2))
	startPos := endPos.Sub(source.Bounds().Size())
	draw.Draw(dest, image.Rect(startPos.X, startPos.Y, endPos.X, endPos.Y), source, image.Point{}, draw.Over)
}
