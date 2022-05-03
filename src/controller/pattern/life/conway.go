package life

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/pixel"
	"image"
	"math/rand"
	"time"
)

type Pattern struct {
	currentBoardCounter uint
	frameDuration       time.Duration
	boards              []*gg.Context
	maxX                int
	maxY                int
	channel             chan image.Image
	closed              chan struct{}
}

func Run(maxX, maxY int, filled float32, frameDuration time.Duration) *Pattern {
	p := &Pattern{
		boards: []*gg.Context{
			gg.NewContext(maxX, maxY),
			gg.NewContext(maxX, maxY),
		},
		frameDuration: frameDuration,
	}

	for _, graphic := range p.boards {
		graphic.Clear()
	}

	p.currentImage().SetColor(color.Crimson)
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			if rand.Float32() <= filled {
				p.currentImage().SetPixel(x, y)
			}
		}
	}

	go p.run()

	return p
}

func (p *Pattern) currentImage() *gg.Context {
	return p.boards[p.currentBoardCounter%2]
}

func (p *Pattern) nextImage() *gg.Context {
	return p.boards[(p.currentBoardCounter+1)%2]
}

// clear the current board
// and advance the counter
func (p *Pattern) flopImage() {
	p.currentImage().Clear()
	p.currentBoardCounter += 1
}

func (p *Pattern) stepPixel(x, y int) {
	numNeighbors := 0
	if x != 0 {
		if !pixel.IsEmpty(p.currentImage(), x-1, y) {
			numNeighbors += 1
		}
	}

	if x != p.maxX {
		if !pixel.IsEmpty(p.currentImage(), x+1, y) {
			numNeighbors += 1
		}
	}

	if y != 0 {
		if !pixel.IsEmpty(p.currentImage(), x, y-1) {
			numNeighbors += 1
		}
	}

	if y != p.maxY {
		if !pixel.IsEmpty(p.currentImage(), x, y+1) {
			numNeighbors += 1
		}
	}

	if numNeighbors == 2 {
		pixel.Copy(p.currentImage(), p.nextImage(), x, y)
	} else if numNeighbors == 3 {
		if pixel.IsEmpty(p.currentImage(), x, y) {
			p.nextImage().SetColor(color.Maroon)

		} else {
			pixel.Copy(p.currentImage(), p.nextImage(), x, y)
		}
	}
}

func (p *Pattern) Stop() {
	p.closed <- struct{}{}
	<-p.closed
}

func (p *Pattern) run() {
	t := time.Tick(p.frameDuration)
	defer func() {
		close(p.channel)
		p.closed <- struct{}{}
	}()
	for {
		for x := 0; x < p.maxX; x++ {
			for y := 0; y < p.maxY; y++ {
				p.stepPixel(x, y)
			}
		}
		p.flopImage()
		select {
		case <-t:
		case <-p.closed:
			return
		}

		select {
		case p.channel <- p.currentImage().Image():
		case <-p.closed:
			return
		}
	}
}
