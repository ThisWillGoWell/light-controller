package graident

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/graphics"
	"image"
	color "image/color"
	"image/draw"
	"math"
	"time"
)

type HandColors struct {
	HourHand   color.Color
	MinuteHand color.Color
	SecondHand color.Color
}

func (h HandColors) toList() []color.Color {
	return []color.Color{
		h.SecondHand, h.MinuteHand, h.HourHand,
	}
}

func CoolBluesColorSet() HandColors {
	return HandColors{
		HourHand:   color.RGBA{R: 51, G: 153, B: 255, A: 255}, // #3399FF
		MinuteHand: color.RGBA{G: 204, B: 153, A: 255},        // #00CC99
		SecondHand: color.RGBA{G: 153, B: 255, A: 255},        // #0099FF
	}
}

func CalmGreensColorSet() HandColors {
	return HandColors{
		HourHand:   color.RGBA{G: 204, B: 153, A: 255},       // #00CC99
		MinuteHand: color.RGBA{R: 51, G: 204, B: 51, A: 255}, // #33CC33
		SecondHand: color.RGBA{R: 102, G: 204, A: 255},       // #66CC00
	}
}

func pastelDelightsColorSet() HandColors {
	return HandColors{
		HourHand:   color.RGBA{R: 204, G: 204, B: 255, A: 255}, // #CCCCFF
		MinuteHand: color.RGBA{R: 255, G: 204, B: 153, A: 255}, // #FFCC99
		SecondHand: color.RGBA{R: 255, G: 204, B: 204, A: 255}, // #FFCCCC
	}
}

func icyCoolColorSet() HandColors {
	return HandColors{
		HourHand:   color.RGBA{R: 153, G: 204, B: 255, A: 255}, // #99CCFF
		MinuteHand: color.RGBA{R: 102, G: 153, B: 204, A: 255}, // #6699CC
		SecondHand: color.RGBA{R: 204, G: 255, B: 255, A: 255}, // #CCFFFF
	}
}

func OceanBreezeColorSet() HandColors {
	return HandColors{
		HourHand:   color.RGBA{R: 102, G: 204, B: 204, A: 255}, // #66CCCC
		MinuteHand: color.RGBA{G: 102, B: 153, A: 255},         // #006699
		SecondHand: color.RGBA{R: 51, G: 204, B: 204, A: 255},  // #33CCCC
	}
}

func TwilightSerenityColorSet() HandColors {
	return HandColors{
		HourHand:   color.RGBA{R: 102, G: 102, B: 153, A: 255}, // #666699
		MinuteHand: color.RGBA{R: 51, G: 51, B: 153, A: 255},   // #333399
		SecondHand: color.RGBA{G: 51, B: 102, A: 255},          // #003366
	}
}

type Config struct {
	Size   int
	radius float64
	Colors HandColors
}

func Clock(c Config) func(img draw.Image, pos image.Point) {
	dc := gg.NewContext(c.Size, c.Size)
	c.radius = float64(c.Size) / 2
	return func(img draw.Image, pos image.Point) {
		drawClock(dc, img, pos, c)
	}
}

func drawClock(dc *gg.Context, dest draw.Image, pos image.Point, c Config) {
	dc.Clear()
	// Get the current time
	now := time.Now()
	hour, min, sec := now.Hour(), now.Minute(), float64(now.Second())

	// Calculate the angles for the hour, minute, and second hands
	hourAngle := (float64(hour%12)*30.0 + float64(min)/2.0) * math.Pi / 180.0
	minuteAngle := (float64(min) * 6.0) * math.Pi / 180.0
	secondAngle := (sec * 6.0) * math.Pi / 180.0

	// Draw the circular gradient
	drawCircularGradient(dc, c.Colors.toList(), c.radius, hourAngle, minuteAngle, secondAngle)

	// Copy the image to the destination
	draw.Draw(dest, image.Rect(pos.X, pos.Y, pos.X+c.Size, pos.Y+c.Size), dc.Image(), image.Point{}, draw.Src)
}
func drawCircularGradient(dc *gg.Context, colors []color.Color, radius float64, angles ...float64) {
	// the center is also radius, radius
	x, y := radius, radius

	totalAngle := 2 * math.Pi // Assuming angles cover a complete circle
	totalSegments := 360      // Total number of segments for a complete circle

	for i := 0; i < len(colors); i++ {
		startColor := graphics.ColorToRGBA(colors[i])
		endColor := graphics.ColorToRGBA(colors[(i+1)%len(colors)])

		// Calculate the angle and segments for this color
		angleStart := angles[i]
		angleEnd := angles[(i+1)%len(angles)]
		angleDiff := angleEnd - angleStart
		if angleDiff < 0 {
			angleDiff += 2 * math.Pi // Adjust for wraparound
		}
		segments := int(angleDiff / totalAngle * float64(totalSegments))

		anglePerSegment := angleDiff / float64(segments)
		for j := 0; j < segments; j++ {
			t := float64(j) / float64(segments)

			// Interpolate the color between the two gradient stops using lerp
			lerpColor := lerpColor(startColor, endColor, t)

			// Calculate the endpoint for the gradient segment
			currentAngle := angleStart + anglePerSegment*float64(j)
			nextAngle := currentAngle + anglePerSegment

			// Draw the gradient segment
			dc.SetRGBA255(int(lerpColor.R), int(lerpColor.G), int(lerpColor.B), int(lerpColor.A))
			dc.MoveTo(x, y)
			dc.LineTo(x+radius*math.Cos(currentAngle), y+radius*math.Sin(currentAngle))
			dc.DrawArc(x, y, radius, float64(currentAngle), float64(nextAngle))
			dc.ClosePath()
			dc.Fill()
		}
	}
}

func lerpColor(startColor, endColor color.RGBA, t float64) color.RGBA {
	// Linear interpolation (lerp) between two colors
	r := uint8(float64(startColor.R)*(1-t) + float64(endColor.R)*t)
	g := uint8(float64(startColor.G)*(1-t) + float64(endColor.G)*t)
	b := uint8(float64(startColor.B)*(1-t) + float64(endColor.B)*t)
	a := uint8(float64(startColor.A)*(1-t) + float64(endColor.A)*t)
	return color.RGBA{r, g, b, a}
}
