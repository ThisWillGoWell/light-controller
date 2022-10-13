package main

import (
	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/fft"
	"math"
)

const sampleRate = 44100

// samples per second
const processingRate = 100.0

func main() {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}

	framesPerBuffer := int(math.Ceil(sampleRate / processingRate))
	defer portaudio.Terminate()
	buffer := make([]float64, framesPerBuffer)
	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buffer), func(in []float32) {
		for i := range buffer {
			buffer[i] = float64(in[i])
		}
		fftOutput := fft.FFTReal(buffer)
		_ = fftOutput
	})
	if err != nil {
		panic(err)
	}
	stream.Start()
	defer stream.Close()

	select {}
}
