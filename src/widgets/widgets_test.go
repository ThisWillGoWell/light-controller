package widgets

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"github.com/thiswillgowell/light-controller/src/piportal/portals"
	"github.com/thiswillgowell/light-controller/src/widgets/clocks/digital"
	"image"
	"image/color"
	"testing"
	"time"
)

func TestWidget(t *testing.T) {

	//go func() {
	//	for {
	//		weather.VerticalFullInfoDraw(WeatherDispaly.Image(), image.Point{})
	//		WeatherDispaly.Update()
	//	}
	//}()

	go func() {
		//drawClock := graident.Clock(graident.Config{
		//	Size: 192,
		//	Colors: graident.HandColors{
		//		HourHand:   colornames.Red,
		//		MinuteHand: colornames.Blue,
		//		SecondHand: colornames.Blue,
		//	},
		//})
		//for {
		//	drawClock(ClockDisplay.Image(), image.Point{X: -1 * 192 / 4, Y: 0})
		//	ClockDisplay.Update()
		//}
		digitalClock := digital.MilliClock(text.Medium, color.White)
		for {

			rightLong := display.NewRotation(portals.BothVerticals, display.Clockwise)

			digitalClock(rightLong.Image(), image.Point{X: 0, Y: 0})
			rightLong.Update()
		}
	}()
	<-time.After(time.Hour)
}
