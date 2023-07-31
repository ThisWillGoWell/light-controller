package digital

import (
	"fmt"
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"image"
	"image/color"
	"image/draw"
	"time"
)

func SecondClock(font text.BitFontType, c color.Color) func(img draw.Image, pos image.Point) {
	return func(img draw.Image, pos image.Point) {
		now := time.Now()
		if now.Hour() < 10 {
			pos.X += text.WriteOnImage(fmt.Sprintf(" %d", now.Hour()), font, c, pos, img).X
		} else {
			pos.X += text.WriteOnImage(fmt.Sprintf("%d", now.Hour()), font, c, pos, img).X
		}
		pos.X += writeTrimmedSemiColon(img, pos, font, c).X
		pos.X += text.WriteOnImage(now.Format("04"), font, c, pos, img).X
		pos.X += writeTrimmedSemiColon(img, pos, font, c).X
		pos.X = text.WriteOnImage(now.Format("05"), font, c, pos, img).X
	}
}

func MilliClock(font text.BitFontType, c color.Color) func(img draw.Image, pos image.Point) {
	return func(img draw.Image, pos image.Point) {
		now := time.Now()
		if now.Hour() < 10 {
			pos.X += text.WriteOnImage(fmt.Sprintf(" %d", now.Hour()), font, c, pos, img).X
		} else {
			pos.X += text.WriteOnImage(fmt.Sprintf("%d", now.Hour()), font, c, pos, img).X
		}
		pos.X += writeTrimmedSemiColon(img, pos, font, c).X
		pos.X += text.WriteOnImage(now.Format("04"), font, c, pos, img).X
		pos.X += writeTrimmedSemiColon(img, pos, font, c).X
		pos.X += text.WriteOnImage(now.Format("05"), font, c, pos, img).X
		pos.X += text.WriteOnImage(now.Format(".000"), font, c, pos, img).X
	}
}

func writeTrimmedSemiColon(img draw.Image, pos image.Point, font text.BitFontType, c color.Color) image.Point {
	trim := font.EmptyPixels(':')
	pos.X -= trim
	return text.WriteOnImage(":", font, c, pos, img).Sub(image.Point{X: trim})
}
