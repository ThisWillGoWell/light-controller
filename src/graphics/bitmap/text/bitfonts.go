package text

import (
	_ "embed"
	"go.uber.org/zap"
	"image"
	"image/color"
	"image/draw"
	"math"
)
import bdf "github.com/zachomedia/go-bdf"

//go:embed fonts/10x20.bdf
var fontFile10x20 []byte

//go:embed fonts/4x6.bdf
var fontFile4x6 []byte

//go:embed fonts/5x7.bdf
var fontFile5x7 []byte

//go:embed fonts/5x8.bdf
var fontFile5x8 []byte

//go:embed fonts/6x10.bdf
var fontFile6x10 []byte

//go:embed fonts/6x12.bdf
var fontFile6x12 []byte

//go:embed fonts/6x13.bdf
var fontFile6x13 []byte

//go:embed fonts/6x13B.bdf
var fontFile6x13B []byte

//go:embed fonts/6x13O.bdf
var fontFile6x13O []byte

//go:embed fonts/6x9.bdf
var fontFile6x9 []byte

//go:embed fonts/7x13.bdf
var fontFile7x13 []byte

//go:embed fonts/7x13B.bdf
var fontFile7x13B []byte

//go:embed fonts/7x13O.bdf
var fontFile7x13O []byte

//go:embed fonts/7x14.bdf
var fontFile7x14 []byte

//go:embed fonts/7x14B.bdf
var fontFile7x14B []byte

//go:embed fonts/8x13.bdf
var fontFile8x13 []byte

//go:embed fonts/8x13B.bdf
var fontFile8x13B []byte

//go:embed fonts/8x13O.bdf
var fontFile8x13O []byte

//go:embed fonts/9x15.bdf
var fontFile9x15 []byte

//go:embed fonts/9x15B.bdf
var fontFile9x15B []byte

//go:embed fonts/9x18.bdf
var fontFile9x18 []byte

//go:embed fonts/9x18B.bdf
var fontFile9x18B []byte

//go:embed fonts/clR6x12.bdf
var fontFileClR6x12 []byte

//go:embed fonts/helvR12.bdf
var fontFileHelvR12 []byte

//go:embed fonts/texgyre-27.bdf
var fontFileTexgyte []byte

type BitFontType int

func (t BitFontType) Font() *font {
	return fonts[t]
}

const (
	Font10x20Type BitFontType = iota
	Font4x6Type   BitFontType = iota
	Font5x7Type   BitFontType = iota
	Font5x8Type   BitFontType = iota
	Font6x10Type  BitFontType = iota
	Font6x12Type  BitFontType = iota
	Font6x13Type  BitFontType = iota
	Font6x13BType BitFontType = iota
	Font6x13OType BitFontType = iota
	Font6x9Type   BitFontType = iota
	Font7x13Type  BitFontType = iota
	Font7x13BType BitFontType = iota
	Font7x13OType BitFontType = iota
	Font7x14Type  BitFontType = iota
	Font7x14BType BitFontType = iota
	Font8x13Type  BitFontType = iota
	Font8x13BType BitFontType = iota
	Font8x13OType BitFontType = iota
	Font9x15Type  BitFontType = iota
	Font9x15BType BitFontType = iota
	Font9x18Type  BitFontType = iota
	Font9x18BType BitFontType = iota
	ClR6x12Type   BitFontType = iota
	HelvR12Type   BitFontType = iota
	TexgyteType   BitFontType = iota
)

var typeToByte = map[BitFontType][]byte{
	Font10x20Type: fontFile10x20,
	Font4x6Type:   fontFile4x6,
	Font5x7Type:   fontFile5x7,
	Font5x8Type:   fontFile5x8,
	Font6x10Type:  fontFile6x10,
	Font6x12Type:  fontFile6x12,
	Font6x13Type:  fontFile6x13,
	Font6x13BType: fontFile6x13B,
	Font6x13OType: fontFile6x13O,
	Font6x9Type:   fontFile6x9,
	Font7x13Type:  fontFile7x13,
	Font7x13BType: fontFile7x13B,
	Font7x13OType: fontFile7x13O,
	Font7x14Type:  fontFile7x14,
	Font7x14BType: fontFile7x14B,
	Font8x13Type:  fontFile8x13,
	Font8x13BType: fontFile8x13B,
	Font8x13OType: fontFile8x13O,
	Font9x15Type:  fontFile9x15,
	Font9x15BType: fontFile9x15B,
	Font9x18Type:  fontFile9x18,
	Font9x18BType: fontFile9x18B,
	ClR6x12Type:   fontFileClR6x12,
	HelvR12Type:   fontFileHelvR12,
	TexgyteType:   fontFileTexgyte,
}

var fonts = map[BitFontType]*font{}

type font struct {
	*bdf.Font
	emptyEdgeSpace map[rune]int
}

func (t BitFontType) Height() int {
	return fonts[t].Size
}

func (t BitFontType) Width() int {
	return fonts[t].XHeight
}

func (t BitFontType) EmptyPixels(chr rune) int {
	value, ok := fonts[t].emptyEdgeSpace[chr]
	if !ok {
		value = calculateEmptySpace(fonts[t].CharMap[chr])
		fonts[t].emptyEdgeSpace[chr] = value
	}
	return value
}

// calculate the min distance from the edge of the image to a column with a non-zero value
func calculateEmptySpace(character *bdf.Character) int {
	// Get the image.Alpha data from the character
	img := character.Alpha

	// Determine the bounds of the image
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Initialize the minimum distance from the left and right edges of the image
	minDistanceLeft := width // Start with the maximum possible value
	minDistanceRight := width

	// Loop through the columns to find the minimum distance from the left and right edges
	for x := 0; x < width; x++ {
		// Check if there's any non-zero (non-empty) pixel in the current column
		isEmptyColumn := true
		for y := 0; y < height; y++ {
			if img.AlphaAt(x, y).A != 0 {
				isEmptyColumn = false
				break
			}
		}

		// If it's an empty column, calculate the distance from the edges
		if !isEmptyColumn {
			distanceFromLeft := x
			distanceFromRight := width - x - 1 // The rightmost column has index (width - 1 - 1)

			// Update the minimum distances if necessary
			if distanceFromLeft < minDistanceLeft {
				minDistanceLeft = distanceFromLeft
			}
			if distanceFromRight < minDistanceRight {
				minDistanceRight = distanceFromRight
			}
		}
	}

	// Return the minimum distance from the edge to a column with a non-zero value
	return int(math.Min(float64(minDistanceLeft), float64(minDistanceRight)))
}

const (
	ExtraSmall = Font4x6Type
	Small      = ClR6x12Type
	Medium     = Font7x13Type
	MediumBold = Font7x13BType
	Large      = Font9x15Type
	LargeBold  = Font9x15BType
	ExtraLarge = Font10x20Type
)

func init() {
	for t, data := range typeToByte {
		bdfFont, err := bdf.Parse(data)
		if err != nil {
			zap.S().Errorw("failed to parse font data", "position", t, zap.Error(err))
		} else {
			fonts[t] = &font{
				Font:           bdfFont,
				emptyEdgeSpace: map[rune]int{},
			}
		}
	}
}

var colors = map[color.Color]*image.RGBA{}

const (
	maxXCharacterSize = 20
	maxYCharacterSize = 20
)

func createFill(c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, maxXCharacterSize, maxXCharacterSize))
	for x := 0; x < maxXCharacterSize; x++ {
		for y := 0; y < maxYCharacterSize; y++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func WriteOnImage(message string, fontType BitFontType, c color.Color, position image.Point, img draw.Image) image.Point {
	startPosition := position
	fillColor := colors[c]
	if fillColor == nil {
		colors[c] = createFill(c)
		fillColor = colors[c]
	}
	font := fonts[fontType]

	if font == nil {
		zap.S().Errorw("font not found", "type", fontType)
		return image.Point{}
	}
	calculateDistance := func() image.Point {
		distance := position.Sub(startPosition)
		distance = distance.Add(image.Point{X: 0, Y: font.Size})
		return distance
	}

	for _, char := range message {
		character := font.CharMap[char]
		if character == nil {
			zap.S().Errorw("character not found in font", "char", char)
			return calculateDistance()
		}
		bounds := character.Alpha.Bounds()
		draw.DrawMask(img, bounds.Add(position), fillColor, image.Point{X: 0, Y: 0}, character.Alpha, image.Point{X: 0, Y: 0}, draw.Src)
		position = position.Add(image.Point{X: bounds.Max.X})
	}

	return calculateDistance()
}
