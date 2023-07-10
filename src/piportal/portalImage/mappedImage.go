package portalImage

import (
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

type mappedRead func(x1, y1 int) (int, int)

type Image struct {
	*image.RGBA
	mappedImage *mappedImage
}

func (m *Image) Image() draw.Image {
	//TODO implement me
	panic("implement me")
}

func (m *Image) Update() {
	//TODO implement me
	panic("implement me")
}

func NewMappedImage(x, y, x2, y2 int) *Image {
	img := image.NewRGBA(image.Rect(0, 0, x, y))
	return &Image{
		RGBA: img,
		mappedImage: &mappedImage{
			img:               img,
			bounds:            image.Rect(0, 0, x2, y2),
			mappingTransforms: initMappingFunc(x2, y2),
		},
	}
}

func initMappingFunc(x, y int) [][]*mappingTransform {
	reads := make([][]*mappingTransform, x)
	for i := range reads {
		reads[i] = make([]*mappingTransform, y)
	}
	return reads
}

func (m *Image) MappedImage() *mappedImage {
	return m.mappedImage
}

func (m *Image) RegisterMapTransform(t mappingTransform) {
	for x := t.startX; x <= t.endX; x++ {
		for y := t.startY; y <= t.endY; y++ {
			m.mappedImage.mappingTransforms[x][y] = &t
		}
	}
}

type mappedImage struct {
	mappingTransforms [][]*mappingTransform
	bounds            image.Rectangle
	img               *image.RGBA
}

func (m *mappedImage) ColorModel() color.Model {
	return m.img.ColorModel()
}

func (m *mappedImage) Bounds() image.Rectangle {
	return m.bounds
}

func (m *mappedImage) At(x, y int) color.Color {
	transform := m.mappingTransforms[x][y]
	if transform != nil {
		x, y = transform.transform(x, y)
	}
	return m.img.At(x, y)
}

type mappingTransform struct {
	mappingFuncs []mappedRead
	panelHeight  int
	panelWidth   int
	startX       int
	startY       int
	endY         int
	endX         int
}

func (t mappingTransform) transform(x, y int) (int, int) {
	for _, f := range t.mappingFuncs {
		x, y = f(x, y)
	}
	return x, y
}

func (t mappingTransform) translation(dx, dy int) mappingTransform {
	t.mappingFuncs = append(t.mappingFuncs, func(x1, y1 int) (int, int) {
		return x1 + dx, y1 + dy
	})
	return t
}

func (t mappingTransform) rotate180() mappingTransform {
	t.mappingFuncs = append(t.mappingFuncs, func(x1, y1 int) (int, int) {
		// zero the cords and rotate
		x1 = t.endX - (x1 - t.startX)
		y1 = t.endY - (y1 - t.startY)
		return x1, y1
	})
	return t
}

func (t mappingTransform) flipY() mappingTransform {
	t.mappingFuncs = append(t.mappingFuncs, func(x1, y1 int) (int, int) {
		return x1, t.endY - (y1 - t.startY)
	})
	return t
}

func (t mappingTransform) rotateCCW() mappingTransform {
	t.mappingFuncs = append(t.mappingFuncs, func(x1, y1 int) (int, int) {
		return y1, x1
	})
	return t
}

func (t mappingTransform) rotateCW() mappingTransform {
	t.mappingFuncs = append(t.mappingFuncs, func(x1, y1 int) (int, int) {
		return y1, t.endX - x1 - 1
	})
	return t
}
