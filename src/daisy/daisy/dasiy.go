package daisy

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/thiswillgowell/light-controller/ratetracker"

	"github.com/sirupsen/logrus"

	"github.com/jacobsa/go-serial/serial"
)

const NumFrequencies = 88
const bufferSize = 2 * NumFrequencies * 4 //float32 * numFrequencies, left and right

type Daisy struct {
	port       io.ReadWriteCloser
	FFTChannel chan [][]float32
	close      chan struct{}
	tracker    *ratetracker.Tracker
}

func Float32fromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func InitDaisy() (*Daisy, error) {
	options := serial.OpenOptions{
		PortName:        port, // populated from file depending on OS
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: bufferSize,
	}

	port, err := serial.Open(options)
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
			for i := 0; i < NumFrequencies*2; i++ {
				channels[i/NumFrequencies][i%NumFrequencies] = Float32fromBytes(readBuffer[i*4 : i*4+4])
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
