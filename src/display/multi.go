package display

import (
	"fmt"
	"strings"
	"sync"

	"github.com/thiswillgowell/light-controller/color"
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
}

func (m *MultiDisplay) Rows() int {
	return m.NumRows
}

func (m *MultiDisplay) Cols() int {
	return m.NumCols
}

func (m *MultiDisplay) Clear() {
	for _, d := range m.Displays {
		d.Clear()
	}
}

func (m *MultiDisplay) GetPixel(row, col int) color.Color {
	displayIndex := 0
	switch m.Arrangement {
	case ArrangementVertical:
		displayIndex = m.DisplayLookup[row]
		return m.Displays[displayIndex].GetPixel(row-m.StartLocations[displayIndex], col)
	case ArrangementHorizontal:
		displayIndex = m.DisplayLookup[col]
		return m.Displays[displayIndex].GetPixel(row, col-m.StartLocations[displayIndex])
	}
	return color.Off
}

func (m *MultiDisplay) SetPixel(row, col int, c color.Color) {
	displayIndex := 0
	switch m.Arrangement {
	case ArrangementVertical:
		displayIndex = m.DisplayLookup[row]
		m.Displays[displayIndex].SetPixel(row-m.StartLocations[displayIndex], col, c)
	case ArrangementHorizontal:
		displayIndex = m.DisplayLookup[col]
		m.Displays[displayIndex].SetPixel(row, col-m.StartLocations[displayIndex], c)
	}
}

func (m *MultiDisplay) Send() error {
	m.wg.Add(len(m.Displays))
	errors := make([]string, len(m.Displays))
	foundErr := false
	for i, d := range m.Displays {
		go func(i int, d Display) {
			if err := d.Send(); err != nil {
				foundErr = true
				errors[i] = err.Error()
			}
			m.wg.Done()
		}(i, d)
	}
	m.wg.Wait()

	if foundErr {
		return fmt.Errorf("failed to write to sub-display errs=", strings.Join(errors, ","))
	}
	return nil
}

func NewMultiDisplay(arrangement MultiDisplayArrangement, displays ...Display) *MultiDisplay {
	md := &MultiDisplay{
		Arrangement: arrangement,
		Displays:    displays,
		wg:          &sync.WaitGroup{},
	}

	count := 0
	enforceSize := 0
	startCount := 0
	for displayIndex, d := range displays {

		switch arrangement {
		case ArrangementVertical:
			count = d.Rows()
			if enforceSize == 0 {
				enforceSize = d.Cols()
			} else if enforceSize != d.Cols() {
				panic("incorrect number of cols, display not square")
			}
		case ArrangementHorizontal:
			count = d.Cols()
			if enforceSize == 0 {
				enforceSize = d.Rows()
			} else if enforceSize != d.Rows() {
				panic("incorrect number of rows, display not square")
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

	return md

}
