package stream

import (
	"log"
	"time"

	/*
		#include <sensormockery/mock.h>
		#include <sensormockery/mock.c>
	*/
	"C"
)

// Stream represents a stream.
type Stream struct {
	ID         int
	WaveType   string
	Sensor     string
	NoiseCoeff float64
	BrokerURL  string
}

const (
	// SawtoothWave is a sawtooth wave.
	SawtoothWave = "sawtooth"
	// RectangularWave is rectangular wave.
	RectangularWave = "rectangular"
	// TriangularWave is triangular wave.
	TriangularWave = "triangular"
	// SineWave is sine wave.
	SineWave = "sine"
)

const (
	// Accelerometer represents the sensor name.
	Accelerometer = "accelerometer"
	// Gyroscope represents the sensor name.
	Gyroscope = "gyroscope"
)

var streams chan *Stream
var currStreams map[int]bool

func init() {
	streams = make(chan *Stream)
	currStreams = make(map[int]bool)
}

// StartStreamListener starts a listener that when triggered creates streams.
func StartStreamListener() {
	for stream := range streams {
		go sendMocks(stream)
	}
}

// StartStream starts a stream by ID.
func StartStream(stream *Stream) {
	streams <- stream
	currStreams[stream.ID] = true
}

// StopStream stops by ID.
func StopStream(id int) {
	currStreams[id] = false
}

func sendMocks(stream *Stream) {
	deltaTime := 0.0

	for {
		if !currStreams[stream.ID] {
			delete(currStreams, stream.ID)
			return
		}

		deltaTime += 0.02
		mockedData := getMockedData(deltaTime, stream)

		log.Printf("value: %f\n", mockedData)
		time.Sleep(time.Millisecond * 20)
	}
}

func waveType(waveType string) uint32 {
	switch waveType {
	case SawtoothWave:
		return C.Sawtooth
	case RectangularWave:
		return C.Rectangular
	case TriangularWave:
		return C.Triangular
	}

	return C.Sine
}

func getMockedData(x float64, stream *Stream) float64 {
	xInC := C.double(x)
	waveTypeInC := waveType(stream.WaveType)
	noiseCoeffInC := C.double(stream.NoiseCoeff)

	var mockedData C.double

	switch stream.Sensor {
	case Accelerometer:
		mockedData = C.accelerometerMock(xInC, waveTypeInC, noiseCoeffInC)
		break
	case Gyroscope:
		mockedData = C.gyroscopeMock(xInC, waveTypeInC, noiseCoeffInC)
		break
	}

	return float64(mockedData)
}
