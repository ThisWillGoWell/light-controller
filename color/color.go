package color

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Color struct {
	H    uint16
	S    uint8
	V    uint8
	R    uint8
	G    uint8
	B    uint8
	VMod int // 10 vmod in one V
}

type Hsv struct {
	H uint16
	S uint8
	V uint8
}

func (c Color) Darken() Color {

	//remove 1  Vmod
	if c.VMod == 0 && c.V != 0 {
		c.VMod = 3
		c.V = c.V - 1
	} else if c.VMod != 0 {
		c.VMod -= 1
	}

	return c.UpdateHsv()
}

func (c Color) SetHue16(hue uint16) Color {
	c.H = hue
	return c.UpdateHsv()
}

func (c Color) SetHue(hue int) Color {
	c.H = uint16(hue)
	return c.UpdateHsv()
}

func (c Color) AddHue(hue int) Color {
	if hue < 0 {
		c.H -= uint16(hue * -1)
	} else {
		c.H += uint16(hue)
	}
	return c.UpdateHsv()
}

func (c Color) UpdateHsv() Color {
	newC := FromHsv(c.H, c.S, c.V)
	newC.VMod = c.VMod
	return newC
}

func FromHsv(h uint16, s, v uint8) Color {
	r, g, b := HsvToRgb(h, s, v)
	return Color{
		H: h,
		S: s,
		V: v,
		R: r,
		G: g,
		B: b,
	}
}

func FromHexString(s string) Color {
	if len(s) != 6 {
		logrus.Error("got hex string with wrong length, len=%d", len(s))
		return Color{}
	}
	r2, err := strconv.ParseInt(s[0:2], 16, 64)
	if err != nil {
		logrus.Errorf("failed to parse hex string err=%v", err)
		return Color{}
	}

	g2, err := strconv.ParseInt(s[2:4], 16, 64)
	if err != nil {
		logrus.Errorf("failed to parse hex string err=%v", err)
		return Color{}
	}

	b2, err := strconv.ParseInt(s[4:6], 16, 64)
	if err != nil {
		logrus.Errorf("failed to parse hex string err=%v", err)
		return Color{}
	}

	r, g, b := uint8(r2), uint8(g2), uint8(b2)
	h, sat, v := RgbToHsv(r, g, b)

	return Color{
		H: h,
		S: sat,
		V: v,
		R: r,
		G: g,
		B: b,
	}
}

func (c Color) RandBetween(c2 Color) Color {
	newC := Color{}
	if c.R < c2.R {
		newC.R = c.R + uint8(rand.Int63n(int64(c2.R-c.R)))
	} else {
		newC.R = c2.R + uint8(rand.Int63n(int64(c.R-c2.R)))
	}

	if c.G < c2.G {
		newC.G = c.G + uint8(rand.Int63n(int64(c2.G-c.G)))
	} else {
		newC.G = c2.G + uint8(rand.Int63n(int64(c.G-c2.G)))
	}

	if c.B < c2.B {
		newC.B = c.B + uint8(rand.Int63n(int64(c2.B-c.B)))
	} else {
		newC.B = c2.B + uint8(rand.Int63n(int64(c.B-c2.B)))
	}

	newC.H, newC.V, newC.V = RgbToHsv(newC.R, newC.G, newC.B)

	return newC
}

// Hsv returns the Hsv [0..360], Saturation and Value [0..1] of the color.
func RgbToHsv(red, green, blue uint8) (uint16, uint8, uint8) {
	R := float64(red) / 255
	G := float64(green) / 255
	B := float64(blue) / 255
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

/// HSV ///
///////////
// From http://en.wikipedia.org/wiki/HSL_and_HSV
// Note that h is in [0..360] and s,v in [0..1]

// Hsv returns the Hsv [0..360], Saturation and Value [0..1] of the color.
func HsvToRgb(hue uint16, sat, value uint8) (uint8, uint8, uint8) {
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

	return uint8((m + r) * 255), uint8((m + g) * 255), uint8((m + b) * 255)
}
