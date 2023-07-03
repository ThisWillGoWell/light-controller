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

const bufferSize = sampleSize * samplesPerFFT

type microphone struct {
	props   audio.Props
	lock    *sync.Mutex
	running bool
	logger  log.Logger

	stream *portaudio.Stream

	FrequencyChannel chan specturm.FrequencySpectrum

	resourcePool sync.Pool
}

func windowFunction(i int, input float64) float64 {
	scalar := 1.0
	if i < windowingSize {
		scalar = float64(i) / float64(windowingSize)
	}
	if i+windowingSize >= bufferSize {
		scalar = float64(i+windowingSize-bufferSize) / float64(windowingSize)
	}
	return input * scalar
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

	inputBuffer := make([]float64, bufferSize)

	currentBufferIndex := 0

	windowedInput := make([]float64, len(inputBuffer))

	var err error
	d.stream, err = portaudio.OpenDefaultStream(1, 0, sampleRate, framesPerBuffer, func(in []float32) {
		frequencies := specturm.NewSpectrum(resourcePool, 1)

		// copy the data into the input buffer
		for _, value := range in {
			inputBuffer[currentBufferIndex] = float64(value)
			currentBufferIndex = (currentBufferIndex + 1) % len(inputBuffer)
		}

		// reorder the input and apply a window
		for inputIndex, value := range inputBuffer {
			i := (inputIndex - currentBufferIndex) % bufferSize
			windowedInput[i] = windowFunction(i, value)
		}

		// preform the FFT on the windowed input
		fftOutput := fft.FFTReal(windowedInput)

		// the first index is dc offset, and the second half of the array is mirrored
		for i := 1; i < bufferSize/2; i++ {
			frequencies.LeftChannel[i] = specturm.FrequencyValue(mag(fftOutput[i]))
			frequencies.RightChannel[i] = specturm.FrequencyValue(mag(fftOutput[i]))
		}

		d.FrequencyChannel <- frequencies

	})
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	if err := d.stream.Start(); err != nil {
		return err
	}
	return nil
}

func mag(c complex128) float64 {
	return math.Sqrt(math.Pow(real(c), 2) + math.Pow(imag(c), 2))
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
