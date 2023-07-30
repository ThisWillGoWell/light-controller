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
		var distance image.Point
		if now.Hour() < 10 {
			distance = text.WriteOnImage(fmt.Sprintf(" %d", now.Hour()), font, c, pos, img)
		} else {
			distance = text.WriteOnImage(fmt.Sprintf("%d", now.Hour()), font, c, pos, img)
		}

		pos.X += distance.X

		writeTrimmedSemiColon := func() {
			trim := font.Font().EmptyPixels(':')
			pos.X -= trim
			distance = text.WriteOnImage(":", font, c, pos, img)
			pos.X += distance.X - trim
		}

		writeTrimmedSemiColon()
		distance = text.WriteOnImage(now.Format("04"), font, c, pos, img)
		pos.X += distance.X
		writeTrimmedSemiColon()
		text.WriteOnImage(now.Format("05"), font, c, pos, img)
	}
}
