package display

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type MultiDisplayArrangement int

const (
	ArrangementHorizontal MultiDisplayArrangement = iota
	ArrangementVertical
)

type MultiDisplay struct {
	NumRows     int
	NumCols     int
	Arrangement MultiDisplayArrangement
	// given a line alongside the arrangement, return what index the display is that owns it
	DisplayLookup []int
	// list of displays
	Displays []Display
	// given a index of dispaly, return where along the arrangemnt it lies
	StartLocations []int
	wg             *sync.WaitGroup
	BoundingBox    image.Rectangle
}

func (m *MultiDisplay) ColorModel() color.Model {
	return color.RGBAModel
}

func (m *MultiDisplay) Bounds() image.Rectangle {
	return m.BoundingBox
}

func (m *MultiDisplay) At(x, y int) color.Color {
	displayIndex := 0
	switch m.Arrangement {
	case ArrangementVertical:
		displayIndex = m.DisplayLookup[y]
		return m.Displays[displayIndex].Image().At(x, y-m.StartLocations[displayIndex])
	case ArrangementHorizontal:
		displayIndex = m.DisplayLookup[x]
		return m.Displays[displayIndex].Image().At(x-m.StartLocations[displayIndex], y)
	}
	return color.RGBA{}
}

func (m *MultiDisplay) Set(x, y int, c color.Color) {
	displayIndex := 0
	switch m.Arrangement {
	case ArrangementVertical:
		displayIndex = m.DisplayLookup[y]
		m.Displays[displayIndex].Image().Set(x, y-m.StartLocations[displayIndex], c)
	case ArrangementHorizontal:
		displayIndex = m.DisplayLookup[x]
		m.Displays[displayIndex].Image().Set(x-m.StartLocations[displayIndex], y, c)
	}
}

func (m *MultiDisplay) Update() {
	m.wg.Add(len(m.Displays))
	for i, d := range m.Displays {
		go func(i int, d Display) {
			d.Update()
			m.wg.Done()
		}(i, d)
	}
	m.wg.Wait()
}

func (m *MultiDisplay) UpdateImage(newImage image.Image) {
	draw.Draw(m, m.BoundingBox, newImage, image.Point{}, draw.Src)
	m.Update()
}

func (m *MultiDisplay) Image() draw.Image {
	return m
}

func NewMultiDisplay(arrangement MultiDisplayArrangement, displays ...Display) *MultiDisplay {
	md := &MultiDisplay{
		Displays: displays,
		wg:       &sync.WaitGroup{},
	}

	count := 0
	enforceSize := 0
	startCount := 0
	md.Arrangement = arrangement
	for displayIndex, d := range displays {
		switch arrangement {
		case ArrangementVertical:
			count = d.Image().Bounds().Max.Y
			if enforceSize == 0 {
				enforceSize = d.Image().Bounds().Max.X
			} else if enforceSize != d.Image().Bounds().Max.X {
				panic("incorrect number of cols, display not a rectangle")
			}
		case ArrangementHorizontal:
			count = d.Image().Bounds().Max.X
			if enforceSize == 0 {
				enforceSize = d.Image().Bounds().Max.Y
			} else if enforceSize != d.Image().Bounds().Max.Y {
				panic("incorrect number of rows, display not a rectangle")
			}
		}
		for i := 0; i < count; i++ {
			md.DisplayLookup = append(md.DisplayLookup, displayIndex)
		}
		md.StartLocations = append(md.StartLocations, startCount)
		startCount += len(md.DisplayLookup)
	}

	switch arrangement {
	case ArrangementHorizontal:
		md.NumRows = enforceSize
		md.NumCols = len(md.DisplayLookup)
	case ArrangementVertical:
		md.NumRows = len(md.DisplayLookup)
		md.NumCols = enforceSize
	}
	md.BoundingBox = image.Rect(0, 0, md.NumCols, md.NumRows)
	return md

}
