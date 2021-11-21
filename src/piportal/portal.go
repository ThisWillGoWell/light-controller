package piportal

import (
	"github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/piportal/tx"
)

type Matrix struct {
	NumCols int
	NumRows int

	BitDepth int

	Address string

	FrameBuffer [][]color.Color
	WriteBuffer []byte
	Conn        *tx.Connection
}

func (m *Matrix) Height() int {
	return m.NumRows
}

func (m *Matrix) Width() int {
	return m.NumCols
}

func (m *Matrix) Clear() {
	for r := range m.FrameBuffer {
		for c := range m.FrameBuffer[r] {
			m.FrameBuffer[r][c] = color.Off
		}
	}
}

func (m *Matrix) GetPixel(row, col int) color.Color {
	return m.FrameBuffer[row][col]
}

func (m *Matrix) SetPixel(row, col int, c color.Color) {
	m.FrameBuffer[row][col] = c
}

func (m *Matrix) Send() error {
	return m.Write()
}

func NewMatrix(address string) (*Matrix, error) {
	m := &Matrix{
		NumRows:  64,
		NumCols:  64,
		Address:  address,
		BitDepth: 8,
	}
	m.WriteBuffer = make([]byte, m.NumRows*m.NumCols*3)

	for r := 0; r < m.NumRows; r++ {
		m.FrameBuffer = append(m.FrameBuffer, make([]color.Color, m.NumCols))
	}

	conn, err := tx.NewClient(address)

	if err != nil {
		return nil, err
	}
	m.Conn = conn
	return m, nil
}

func (m *Matrix) ForEachAndUpdate(each func(r, c int) color.Color) error {
	for r := 0; r < m.NumRows; r++ {
		for c := 0; c < m.NumCols; c++ {
			m.FrameBuffer[r][c] = each(r, c)
		}
	}
	return m.Write()
}

func (m *Matrix) Write() error {
	err := m.Conn.WriteFrame(m.FrameBuffer)
	return err
}

func (m *Matrix) SetColor() {
	m.ForEachAndUpdate(func(r, c int) color.Color {
		return color.Deeppink
	})
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Pattern struct {
	Row           int
	Col           int
	Direction     Direction
	CurrentOffset int
}

func (p *Pattern) Init(m *Matrix) error {
	return m.ForEachAndUpdate(func(r, c int) color.Color {
		return color.Color{}
	})
}

func (p *Pattern) Draw(m *Matrix) {
	m.FrameBuffer[p.Row][p.Col] = color.Color{}
}
