package piportal

import (
	"image"
	"image/draw"

	"github.com/thiswillgowell/light-controller/src/piportal/tx"
)

type Matrix struct {
	image *image.RGBA
	Conn  *tx.Connection
}

func (m *Matrix) Image() draw.Image {
	return m.image
}

func (m *Matrix) UpdateImage(src image.Image) {
	draw.Draw(m.image, m.image.Bounds(), src, image.Point{}, draw.Src)
	m.Update()
}

func NewMatrix(address string) (*Matrix, error) {
	m := &Matrix{
		image: image.NewRGBA(image.Rect(0, 0, 64, 96)),
	}
	conn, err := tx.NewTCPClient(address)

	if err != nil {
		return nil, err
	}
	m.Conn = conn
	return m, nil
}

func (m *Matrix) Update() {
	m.Conn.WriteFrame(m.image)
}
