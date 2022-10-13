package fire

import (
	"github.com/fogleman/gg"
	color2 "github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/graphics"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math/rand"
)

type FlamePattern struct {
	d        display.Display
	dc       *gg.Context
	data     []int
	size     image.Point
	fps      float64
	counter  int
	wind     float64
	strength float64
	colors   []color.Color
}

func NewFlamePattern(d display.Display, fps float64) *FlamePattern {
	return &FlamePattern{
		d:        d,
		dc:       gg.NewContextForImage(d.Image()),
		size:     d.Image().Bounds().Max,
		fps:      fps,
		strength: strMin,
		wind:     windMin,
		colors:   append([]color.Color{color2.Off}, graphics.LinerGradient(colornames.Orangered, colornames.Darkorange, 31)...),
		data:     make([]int, (d.Image().Bounds().Max.Y+1)*(d.Image().Bounds().Max.X+1)),
	}
}

var (
	timeFire0 = 20.0
	timeFire1 = 0.5

	windMin = 0.0
	windMax = 3.0

	strMin = 2.0
	strMax = 4.0
)

func setVal(frame FlamePattern, x int, y int, val int) FlamePattern {
	if x < 0 || x >= frame.size.X {
		return frame
	}
	if y < 0 || y >= frame.size.Y {
		return frame
	}

	frame.data[y*frame.size.X+x] = val
	return frame
}

func (flame *FlamePattern) step() {
	for i := 0; i < flame.size.X; i++ {
		v := 0
		if flame.counter < int(timeFire0*flame.fps) {
			v = 31
		}
		flame.data[(flame.size.Y-1)*flame.size.X+i] = v
	}
	//flame.drawFlame()

	for x := 0; x < flame.size.X; x++ {
		for y := 1; y < flame.size.Y; y++ {
			// NOTE: fire decay
			v := flame.data[(y*flame.size.X)+x]
			if v < 0 {
				continue
			}

			v -= int(rand.Float64() * 2)
			if v < 0 {
				v = 0
			}
			flame.data[y*flame.size.X+x] = v

			// NOTE: fire spread
			y2 := y - int(rand.Float64()*flame.strength)
			if y2 < 0 {
				continue
			}

			spread := flame.strength * 2
			x2 := x + int(rand.Float64()*spread-(spread/2)+flame.wind)
			if x2 < 0 || x2 >= flame.size.X {
				continue
			}
			flame.data[y2*flame.size.X+x2] = v
		}
	}

	flame.counter = (flame.counter + 1) % int((timeFire0+timeFire1)*flame.fps)

	if flame.counter == 0 {
		flame.wind++
		if flame.wind > windMax {
			flame.wind = windMin
		}

		if flame.wind == 0 {
			flame.strength++
			if flame.strength > strMax {
				flame.strength = strMin
			}
		}
	}
	flame.drawFlame()
}

func (flame *FlamePattern) drawFlame() {
	flame.dc.SetColor(color.RGBA{})
	flame.dc.Clear()
	for x := 0; x < flame.size.X; x++ {
		for y := 0; y < flame.size.Y; y++ {
			v := flame.data[(y*flame.size.X)+x]
			flame.dc.SetColor(flame.colors[v])
			flame.dc.SetPixel(x, y)
		}
	}
	display.DrawAndUpdate(flame.d, flame.dc.Image())
}
