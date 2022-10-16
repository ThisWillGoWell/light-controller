package microphone

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/fft"
	"github.com/thiswillgowell/light-controller/log"
	"github.com/thiswillgowell/light-controller/src/audio"
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"math"
	"sync"
)

var Microphone = &microphone{
	lock: &sync.Mutex{},
}

const sampleRate = 96000

// samples per second processed
// used to calculate the block size which should be a power of 2
const idealProcessingRate = 100.0

const sampleSize = 512

const samplesPerFFT = 16

const historicalValues = 16

const windowingSize = 32

type microphone struct {
	props   audio.Props
	lock    *sync.Mutex
	running bool
	logger  log.Logger

	stream *portaudio.Stream

	frequencyChannel chan specturm.FrequencySpectrum

	resourcePool sync.Pool
}

func windowFunction(i int, input float64) float64 {
	if i >= windowingSize && i < ()windowingSize
}

func (d *microphone) Start() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.running {
		d.logger.Warn("starting microphone while its running")
		return nil
	}
	if err := portaudio.Initialize(); err != nil {
		return err
	}

	framesPerBuffer := int(math.Pow(2, math.Ceil(math.Log(sampleRate/idealProcessingRate))))
	resourcePool := specturm.NewSamplePool(framesPerBuffer)

	// append values as they are read

	// pass each completed buffer to a low pass filter

	// window over each buffer as its completed

	// pass values into fft

	fftBuffers := make([][]float32, historicalValues)
	for i := range fftBuffers {
		fftBuffers[i] = make([]float32, sampleSize*samplesPerFFT)
	}

	inputBuffer := make([]float64, sampleSize*samplesPerFFT)

	currentBufferIndex := 0

	currentIndex := 0

	processingBuffer := make([]float64, framesPerBuffer)
	var err error
	d.stream, err = portaudio.OpenDefaultStream(1, 0, sampleRate, framesPerBuffer, func(in []float32) {
		frequencies := specturm.NewSpectrum(resourcePool, 1)

		// copy the data into the input buffer
		for _, value := range in {
			inputBuffer[currentBufferIndex] = float64(value)
			currentBufferIndex = (currentBufferIndex + 1) % len(inputBuffer)
		}

		// preform the FFT
		fftOutput := fft.FFTReal(append(inputBuffer[currentBufferIndex:], inputBuffer[0:currentBufferIndex]...))

		// copy the fft result to the result buffer and apply windowing function

		fftOutput := fft.FFTReal(processingBuffer)
		for i := range fftOutput {
			frequencies.LeftChannel[i] = specturm.FrequencyValue(real(fftOutput[i]))
			frequencies.RightChannel[i] = specturm.FrequencyValue(real(fftOutput[i]))
		}

		d.frequencyChannel <- frequencies
	})
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	if err := d.stream.Start(); err != nil {
		return err
	}
	return nil
}

func (d *microphone) Stop() {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !d.running {
		return
	}
	if err := portaudio.Terminate(); err != nil {
		d.logger.Error("failed to terminate port audio resources", err)
	}
	d.running = false
}
