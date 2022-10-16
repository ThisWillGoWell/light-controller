package specturm

import "sync"

type Format uint64

type FrequencyValue float64

const (
	MonoFormat Format = iota
	StereoFormat
)

type Channel uint64

const (
	LeftChannel Channel = iota
	RightChannel
)

type Props struct {
	Format Format
	// Frequency step per index
	FrequencyDelta int
}

type FrequencySpectrum struct {
	LeftChannel  []FrequencyValue
	RightChannel []FrequencyValue

	Format Format
	pool   *sync.Pool

	MaxValue FrequencyValue
}

func (f FrequencySpectrum) Free() {
	f.pool.Put(f.LeftChannel)
	f.pool.Put(f.RightChannel)
}

func NewSamplePool(bufferSize int) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return make([]FrequencyValue, bufferSize)
		},
	}
}

func NewSpectrum(pool *sync.Pool, maxValue FrequencyValue) FrequencySpectrum {
	f := FrequencySpectrum{
		pool:     pool,
		MaxValue: maxValue,
	}

	f.RightChannel = pool.Get().([]FrequencyValue)
	f.LeftChannel = pool.Get().([]FrequencyValue)

	return f
}

func (f FrequencySpectrum) Bin(channel Channel, binningInput BinInput, removeDeadZone bool) []int {
	var input []FrequencyValue
	switch channel {
	case RightChannel:
		input = f.RightChannel
	case LeftChannel:
		input = f.LeftChannel
	}
	return BinHeight(input, binningInput, f.MaxValue)
}
