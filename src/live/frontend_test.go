package live

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/display"
	"image"
	"image/color"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	d := display.NewSubscription(display.NewRGBA(100, 100))
	go RunServer(d)

	img := image.NewRGBA(d.Image().Bounds())
	dc := gg.NewContextForRGBA(img)
	maxX, maxY := float64(d.Image().Bounds().Max.X), float64(d.Image().Bounds().Max.Y)
	for x := 0; x < d.Image().Bounds().Max.X; x++ {
		dc.Clear()
		grad := gg.NewRadialGradient(maxX/2, maxY/2, 10, maxX/2, maxY/2, 80)
		grad.AddColorStop(0, color.RGBA{10, 50, 0, 255})
		grad.AddColorStop(1, color.RGBA{60, 0, 40, 255})

		dc.SetFillStyle(grad)
		dc.DrawRectangle(1, 1, maxX-2, maxY-2)
		dc.Fill()

		dc.SetColor(color.White)
		dc.DrawCircle(float64(x), float64(x)/float64(d.Image().Bounds().Max.X)*float64(d.Image().Bounds().Max.Y), 10)
		dc.Stroke()
		dc.DrawCircle(float64(x), 48, 80)
		dc.Stroke()
		dc.DrawLine(0, 0, float64(d.Image().Bounds().Max.X), float64(d.Image().Bounds().Max.Y))
		dc.Stroke()
		display.DrawAndUpdate(d, dc.Image())
		<-time.After(time.Second)
	}

}
