package icons

import (
	"bytes"
	_ "embed"
	"git.sr.ht/~sbinet/gg"
	"github.com/thiswillgowell/light-controller/src/graphics"
	"image"
	"image/png"
)

//go:embed 01d@2x.png
var oneD []byte

//go:embed 01n@2x.png
var oneN []byte

//go:embed 02d@2x.png
var twoD []byte

//go:embed 02n@2x.png
var twoN []byte

//go:embed 03d@2x.png
var threeD []byte

//go:embed 03n@2x.png
var threeN []byte

//go:embed 04d@2x.png
var fourD []byte

//go:embed 04n@2x.png
var fourN []byte

//go:embed 10d@2x.png
var tenD []byte

//go:embed 10n@2x.png
var tenN []byte

//go:embed 11d@2x.png
var elevenD []byte

//go:embed 11n@2x.png
var elevenN []byte

//go:embed 13d@2x.png
var thirteenD []byte

//go:embed 13n@2x.png
var thirteenN []byte

//go:embed 50d@2x.png
var fiftyD []byte

//go:embed 50n@2x.png
var fiftyN []byte

var codeToFile = map[string][]byte{
	"01d": oneD,
	"01n": oneN,
	"02d": twoD,
	"02n": twoN,
	"03d": threeD,
	"03n": threeN,
	"04d": fourD,
	"04n": fourN,
	"10d": tenD,
	"10n": tenN,
	"11d": elevenD,
	"11n": elevenN,
	"13d": thirteenD,
	"13n": thirteenN,
	"50d": fiftyD,
	"50n": fiftyN,
}

// default to rainy image if an icon is not found
const defaultImage = "09n"

var codeToPng = map[string]image.Image{}

func init() {
	for code, value := range codeToFile {
		var err error
		codeToPng[code], err = png.Decode(bytes.NewReader(value))
		if err != nil {
			panic("invalid PNG file in icons: " + err.Error())
		}
	}
}

const trim = 17

// icons are 100x100 by default
const trimmedSize = 100.0 - (trim * 2)

func GetWeatherIcon(code string, size int) image.Image {
	icon := codeToPng[code]
	graphics.SaveImage("test-icon", icon)
	dc := gg.NewContext(size, size)
	dc.Scale(float64(size)/trimmedSize, float64(size)/trimmedSize)
	dc.DrawImage(icon, -1*trim, -1*trim)
	graphics.SaveImage("test-icon-2", dc.Image())
	return dc.Image()
}
