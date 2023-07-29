package graident

import (
	"math"
	"time"
)

type GraidentClockApp struct {
}

func drawClock(dc *gg.Context) {
	// Get the current time
	now := time.Now()
	hour, min, sec := now.Hour(), now.Minute(), now.Second()

	// Calculate the angles for the hour, minute, and second hands
	hourAngle := (float64(hour%12)*30.0 + float64(min)/2.0) * math.Pi / 180.0
	minuteAngle := (float64(min) * 6.0) * math.Pi / 180.0
	secondAngle := (float64(sec) * 6.0) * math.Pi / 180.0

	// Calculate the center of the matrix
	centerX, centerY := Width/2, Height/2

	// Draw the circular gradient
	drawCircularGradient(dc, centerX, centerY, Radius, hourAngle, minuteAngle, secondAngle)
}

func drawCircularGradient(dc *gg.Context, x, y int, radius float64, angles ...float64) {
	// Colors for the gradient stops (hour, minute, second)
	colors := []color.RGBA{
		{255, 0, 0, 255}, // Red for hour hand
		{0, 255, 0, 255}, // Green for minute hand
		{0, 0, 255, 255}, // Blue for second hand
	}

	// Draw the circular gradient
	for i, angle := range angles {
		// Calculate the gradient color at this angle
		r, g, b, a := colors[i].RGBA()
		gradientColor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

		// Calculate the endpoint for the gradient segment
		endX := int(float64(x) + radius*math.Sin(angle))
		endY := int(float64(y) - radius*math.Cos(angle))

		// Draw the gradient segment
		dc.SetRGB255(int(gradientColor.R), int(gradientColor.G), int(gradientColor.B))
		dc.DrawLine(float64(x), float64(y), float64(endX), float64(endY))
		dc.Stroke()
	}
}

