package stream

import (
	"log"
	"time"

	/*
		#include <sensormockery/mock.h>
		#include <sensormockery/mock.c>
	*/
	"C"

	"github.com/streadway/amqp"
)
import (
	"encoding/json"
	"fmt"
)

// Stream represents a stream.
type Stream struct {
	ID         int
	WaveType   string
	Sensor     string
	NoiseCoeff float64
	BrokerURL  string
}

// DataPoint represents data point.
type DataPoint struct {
	Data float64 `json:"data_point"`
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

	conn, channel, queue := initBroker(stream)
	defer conn.Close()
	defer channel.Close()

	for {
		if !currStreams[stream.ID] {
			delete(currStreams, stream.ID)
			return
		}

		deltaTime += 0.2

		dataPoint := &DataPoint{
			Data: getMockedData(deltaTime, stream),
		}
		sendMessageToBroker(channel, queue, dataPoint)

		time.Sleep(time.Millisecond * 200)
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

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func initBroker(stream *Stream) (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial(stream.BrokerURL)
	handleError(err, "Can't connect to AMQP")

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	// create (if not existent) the stream queue
	amqpQueue, err := amqpChannel.QueueDeclare(fmt.Sprintf("stream-%d", stream.ID), false, true, false, false, nil)
	handleError(err, "Could not declare `tasks` queue")

	// set up stream queue
	err = amqpChannel.Qos(1, 0, true)
	handleError(err, "Could not configure QoS")

	return conn, amqpChannel, &amqpQueue
}

func sendMessageToBroker(channel *amqp.Channel, queue *amqp.Queue, dataPoint *DataPoint) {
	body, err := json.Marshal(dataPoint)
	handleError(err, "Failed to marshal data")

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	handleError(err, "Failed to publish a message")
}
