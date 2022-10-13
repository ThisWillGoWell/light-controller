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

var Microphone = &microphone{}

const sampleRate = 44100

// samples per second
const processingRate = 100.0

type microphone struct {
	props   audio.Props
	lock    *sync.Mutex
	running bool
	logger  log.Logger

	stream *portaudio.Stream

	frequencyChannel chan specturm.FrequencySpectrum

	resourcePool sync.Pool
}

func (d *microphone) Start() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.running {
		d.logger.Warn("starting microphone while its running")
		return nil
	}
	framesPerBuffer := int(math.Ceil(sampleRate / processingRate))
	resourcePool := specturm.NewSamplePool(framesPerBuffer)

	processingBuffer := make([]float64, framesPerBuffer)
	var err error
	d.stream, err = portaudio.OpenDefaultStream(1, 0, sampleRate, framesPerBuffer, func(in []float32) {
		frequencies := specturm.NewSpectrum(specturm.MonoFormat, resourcePool)
		for i := range in {
			processingBuffer[i] = float64(in[i])
		}
		fftOutput := fft.FFTReal(processingBuffer)

		for i := range fftOutput {
			frequencies.MonoChannel[i] = specturm.FrequencyValue(real(fftOutput[i]))
		}
		d.frequencyChannel <- frequencies
	})
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
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
