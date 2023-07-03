package main

import (
	"fmt"
	"github.com/thiswillgowell/light-controller/src/audio/inputs/microphone"
)

const sampleRate = 44100

// samples per second
const processingRate = 100.0

func main() {

	mic := microphone.Microphone
	if err := mic.Start(); err != nil {
		panic(err)
	}
	for freq := range mic.FrequencyChannel {
		fmt.Printf("%v\n", freq.RightChannel)
	}

}
