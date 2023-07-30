package widgets

import (
	"github.com/thiswillgowell/light-controller/src/piportal/portals"
	"github.com/thiswillgowell/light-controller/src/widgets/clocks/graident"
	"github.com/thiswillgowell/light-controller/src/widgets/weather"
	"golang.org/x/image/colornames"
	"image"
	"testing"
	"time"
)

func TestWidget(t *testing.T) {

	WeatherDispaly := portals.LeftVertical
	ClockDisplay := portals.RightVertical

	go func() {
		for {
			weather.VerticalFullInfoDraw(WeatherDispaly.Image(), image.Point{})
			WeatherDispaly.Update()
		}
	}()

	go func() {
		drawClock := graident.Clock(graident.Config{
			Size: 192,
			Colors: graident.HandColors{
				HourHand:   colornames.Red,
				MinuteHand: colornames.Blue,
				SecondHand: colornames.Blue,
			},
		})
		for {
			drawClock(ClockDisplay.Image(), image.Point{X: -1 * 192 / 4, Y: 0})
			ClockDisplay.Update()
		}
		//digitalClock := digital.SecondClock(text.Large, color.White)
		//for {
		//	digitalClock(ClockDisplay.Image(), image.Point{X: 0, Y: 0})
		//	ClockDisplay.Update()
		//}
	}()
	<-time.After(time.Minute)
}
