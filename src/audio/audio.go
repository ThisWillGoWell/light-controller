package audio

import (
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"time"
)

type Device interface {
	Start() error
	Stop()
	SpectrumChannel() chan specturm.FrequencySpectrum
	Props() Props
}

// Props of a device
type Props struct {
	// number of channels coming from the device
	NumChannels int
	// delta in between reads
	Rate time.Duration
}

type Manager struct {
	ActiveDevice      Device
	RegisteredDevices map[string]Device
}

//func StartDevice(device Device, processFunctions) {
//
//}
