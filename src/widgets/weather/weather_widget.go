package weather

import (
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"image"
	"image/draw"
)

func VerticalFullInfoDraw(img draw.Image, pos image.Point) {
	CurrentTemp(img, pos)
	graphLength := 48
	Next90MinTempGraph(img, image.Point{X: 3, Y: pos.Y + 30}, 10, graphLength)
	TwoDayHighAndLow(img, image.Point{X: graphLength + 6, Y: pos.Y + 29}, text.ExtraSmall)
	// start pos of last item + height of item + offset
	pos.Y += 30 + 10 + 5
	ThreeDayForecast(img, image.Point{X: pos.X + 6, Y: pos.Y}, 5)
}
