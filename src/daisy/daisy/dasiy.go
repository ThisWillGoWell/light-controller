package daisy

import (
	"encoding/binary"
	"fmt"
	"github.com/thiswillgowell/light-controller/ratetracker"
	"io"
	"math"

	"github.com/sirupsen/logrus"

	"github.com/tarm/serial"
)

const NumFrequencies = 88
const bufferSize = 2 * NumFrequencies * 4 //float32 * numFrequencies, left and right

type Daisy struct {
	port       io.ReadWriteCloser
	FFTChannel chan [][]float32
	close      chan struct{}
	tracker    *ratetracker.Tracker
}

func (d *Daisy) NextFFTValues() [][]float32 {
	return <-d.FFTChannel
}

func float32fromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Init() (*Daisy, error) {
	options := &serial.Config{
		Name:     port, // populated from file depending on OS
		Baud:     115200,
		Size:     8,
		StopBits: 1,
	}

	port, err := serial.OpenPort(options)
	if err != nil {
		return nil, fmt.Errorf("serial.Open: %v", err)
	}
	fftChannel := make(chan [][]float32, 3)
	closeChan := make(chan struct{}, 1)

	d := &Daisy{
		port:       port,
		FFTChannel: fftChannel,
		close:      closeChan,
		tracker:    ratetracker.NewTracker("daisy"),
	}

	// start the reader
	readBuffer := make([]byte, bufferSize)
	go func() {
		run := true
		for run {
			//<-time.After(time.Millisecond * 100)
			d.tracker.Track()
			//fmt.Println("write")
			if _, err := d.port.Write([]byte{0x00}); err != nil {
				logrus.Errorf("failed to write start buffer err=%v", err)
				run = false
				continue
			}

			// read left and right
			channels := make([][]float32, 2)
			channels[0] = make([]float32, NumFrequencies)
			channels[1] = make([]float32, NumFrequencies)
			if _, err := d.port.Read(readBuffer); err != nil {
				logrus.Errorf("failed to read err=%v", err)
				run = false
				continue
			}
			//fmt.Printf("%s\n", hex.EncodeToString(readBuffer))

			for i := 0; i < NumFrequencies*2; i++ {
				channels[i/NumFrequencies][i%NumFrequencies] = float32fromBytes(readBuffer[i*4 : i*4+4])
			}
			// first check if we can close
			select {
			case <-d.close:
				run = false
				continue
			default:
			}

			// either send next float or react to close
			select {
			case d.FFTChannel <- channels:
			case <-d.close:
				run = false
			}
		}
		close(d.FFTChannel)
		if err := d.port.Close(); err != nil {
			logrus.Errorf("err while closing connection ")
		}
		logrus.Info("Closeting Daisy Connection")

	}()

	return d, nil
}
