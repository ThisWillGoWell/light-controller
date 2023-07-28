package draw

import (
	"image/color"
	"image/draw"
)

// DrawLine draws a line on the provided draw context with a specified width using Wu's algorithm.
// drawLine draws a line on the provided draw context with a specified width using Wu's algorithm.
func DrawLine(dc draw.Image, c color.Color, width, x1, y1, x2, y2 int) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	if dx == 0 && dy == 0 {
		// Single point, draw as a single pixel.
		dc.Set(x1, y1, c)
		return
	}

	var swap bool
	if dy > dx {
		// Steeper line, swap x and y axes to make it shallow.
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
		swap = true
	}

	if x2 < x1 {
		// Draw from left to right.
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	gradient := float64(dy) / float64(dx)
	y := float64(y1)

	for x := x1; x <= x2; x++ {
		for w := -width / 2; w < width/2; w++ {
			intensity := int(255 - (y-float64(int(y)))*255)

			if swap {
				dc.Set(int(y)+w, x, colorWithAlpha(c, uint8(intensity)))
			} else {
				dc.Set(x, int(y)+w, colorWithAlpha(c, uint8(intensity)))
			}
		}

		y += gradient
	}
}

// colorWithAlpha returns the color with the specified alpha value.
func colorWithAlpha(c color.Color, alpha uint8) color.RGBA {
	r, g, b, _ := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
}

// DrawGradientRect draws a rectangle with a linear gradient fill on the provided draw context.
func DrawGradientRect(dc draw.Image, c1, c2 color.Color, x1, y1, x2, y2 int) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			ratio := float64(x-x1) / float64(x2-x1)
			r, g, b, a := gradientColor(c1, c2, ratio)
			dc.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}
}

// gradientColor calculates the interpolated color between two colors based on a given ratio.
func gradientColor(c1, c2 color.Color, ratio float64) (r, g, b, a uint8) {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r = uint8(float64(r1)*(1-ratio) + float64(r2)*ratio)
	g = uint8(float64(g1)*(1-ratio) + float64(g2)*ratio)
	b = uint8(float64(b1)*(1-ratio) + float64(b2)*ratio)
	a = uint8(float64(a1)*(1-ratio) + float64(a2)*ratio)

	return r, g, b, a
}

// DrawRect draws a filled rectangle on the provided draw context.
func DrawRect(dc draw.Image, c color.Color, x1, y1, x2, y2 int) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			dc.Set(x, y, c)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
