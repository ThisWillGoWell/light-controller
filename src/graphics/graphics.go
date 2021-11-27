package graphics

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

// Hsv returns the Hsv [0..360], Saturation and Value [0..1] of the color_2.
func ColorToHsva(c color.Color) (uint16, uint8, uint8, uint8) {
	r, g, b, a := c.RGBA()
	R := float64(r) / math.MaxInt32
	G := float64(g) / math.MaxInt32
	B := float64(b) / math.MaxInt32
	var h, s, v float64

	min := math.Min(math.Min(R, G), B)
	v = math.Max(math.Max(R, G), B)
	C := v - min

	s = 0.0
	if v != 0.0 {
		s = C / v
	}

	h = 0.0 // We use 0 instead of undefined as in wp.
	if min != v {
		if v == R {
			h = math.Mod((G-B)/C, 6.0)
		}
		if v == G {
			h = (B-R)/C + 2.0
		}
		if v == B {
			h = (R-G)/C + 4.0
		}
		h *= 60.0
		if h < 0.0 {
			h += 360.0
		}
	}
	return uint16(h * 65535 / 360), uint8(s * 255), uint8(v * 255), uint8(a >> 24)
}

/// HSV ///
///////////
// From http://en.wikipedia.org/wiki/HSL_and_HSV
// Note that h is in [0..360] and s,v in [0..1]

// Hsv returns the Hsv [0..360], Saturation and Value [0..1] of the color_2.
func HsvaToColor(hue uint16, sat, value, alpha uint8) color.RGBA {
	H := float64(hue) / 65535 * 360
	S := float64(sat) / 255
	V := float64(value) / 255
	Hp := H / 60.0
	C := V * S
	X := C * (1.0 - math.Abs(math.Mod(Hp, 2.0)-1.0))

	m := V - C
	r, g, b := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= Hp && Hp < 1.0:
		r = C
		g = X
	case 1.0 <= Hp && Hp < 2.0:
		r = X
		g = C
	case 2.0 <= Hp && Hp < 3.0:
		g = C
		b = X
	case 3.0 <= Hp && Hp < 4.0:
		g = X
		b = C
	case 4.0 <= Hp && Hp < 5.0:
		r = X
		b = C
	case 5.0 <= Hp && Hp < 6.0:
		r = C
		b = X
	}
	return color.RGBA{
		R: uint8((m + r) * 255),
		G: uint8((m + g) * 255),
		B: uint8((m + b) * 255),
		A: alpha,
	}
}

func Darken(c color.Color, amount float64) color.Color {
	h, s, v, a := ColorToHsva(c)
	v = uint8(float64(v) / 255 * (1 - amount))
	return HsvaToColor(h, s, v, a)
}

func LinerGradient(startColor, endColor color.Color, numSteps int) []color.Color {
	output := make([]color.Color, numSteps)
	colors := gg.NewLinearGradient(0, 0, float64(numSteps), 0)
	colors.AddColorStop(0, startColor)
	colors.AddColorStop(1, endColor)

	for i := range output {
		output[i] = colors.ColorAt(i, 0)
	}
	return output
}

func HsvToRgb(hue uint16, sat, value uint8) color.Color {
	H := float64(hue) / 65535 * 360
	S := float64(sat) / 255
	V := float64(value) / 255
	Hp := H / 60.0
	C := V * S
	X := C * (1.0 - math.Abs(math.Mod(Hp, 2.0)-1.0))

	m := V - C
	r, g, b := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= Hp && Hp < 1.0:
		r = C
		g = X
	case 1.0 <= Hp && Hp < 2.0:
		r = X
		g = C
	case 2.0 <= Hp && Hp < 3.0:
		g = C
		b = X
	case 3.0 <= Hp && Hp < 4.0:
		g = X
		b = C
	case 4.0 <= Hp && Hp < 5.0:
		r = X
		b = C
	case 5.0 <= Hp && Hp < 6.0:
		r = C
		b = X
	}
	return color.RGBA{
		R: uint8((m + r) * 255),
		G: uint8((m + g) * 255),
		B: uint8((m + b) * 255),
	}
}

func ColorToHsv(inputColor color.Color) (uint16, uint8, uint8) {
	r, g, b, _ := inputColor.RGBA()

	R := float64(r) / math.MaxInt32
	G := float64(g) / math.MaxInt32
	B := float64(b) / math.MaxInt32
	var h, s, v float64

	min := math.Min(math.Min(R, G), B)
	v = math.Max(math.Max(R, G), B)
	C := v - min

	s = 0.0
	if v != 0.0 {
		s = C / v
	}

	h = 0.0 // We use 0 instead of undefined as in wp.
	if min != v {
		if v == R {
			h = math.Mod((G-B)/C, 6.0)
		}
		if v == G {
			h = (B-R)/C + 2.0
		}
		if v == B {
			h = (R-G)/C + 4.0
		}
		h *= 60.0
		if h < 0.0 {
			h += 360.0
		}
	}
	return uint16(h * 65535 / 360), uint8(s * 255), uint8(v * 255)
}
