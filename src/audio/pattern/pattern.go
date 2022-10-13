package pattern

import (
	"github.com/thiswillgowell/light-controller/src/audio"
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"github.com/thiswillgowell/light-controller/src/controller/pattern"
	"sync"
)

type AudioPattern struct {
	patterns []pattern.MusicPattern
	device   audio.Device
}

func RunAudioPattern(device audio.Device, patterns []pattern.MusicPattern, killChannel chan struct{}) error {
	if err := device.Start(); err != nil {
		return err
	}
	defer device.Stop()
	spectrumChannel := device.SpectrumChannel()

	wg := &sync.WaitGroup{}
	for {
		select {
		case <-killChannel:
			return nil
		default:
		}
		select {
		case values := <-spectrumChannel:
			processSpectrum(values, patterns, wg)
		case <-killChannel:
			return nil
		}
	}
}

func processSpectrum(values specturm.FrequencySpectrum, patterns []pattern.MusicPattern, wg *sync.WaitGroup) {
	wg.Add(len(patterns))
	for _, p := range patterns {
		go func(p pattern.MusicPattern) {
			go p.ProcessSpectrum(values)
		}(p)
	}
	wg.Wait()
}
