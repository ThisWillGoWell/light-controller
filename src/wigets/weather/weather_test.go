package weather

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"image"
	"image/color"
	"testing"
)

func TestDrawOnMatrix(t *testing.T) {
	//p := portals.BothVerticals
	p := display.NewTestRGBA(64, 96, "weather")
	DrawCurrentTemp(p.Image(), text.ExtraLarge, image.Point{0, 0}, color.White)
	p.Update()

}

func TestCurrentWeather(t *testing.T) {
	d := display.NewTestRGBA(64, 96, "current")
	CurrentWeather(d.Image(), image.Point{})
	d.Update()
}
