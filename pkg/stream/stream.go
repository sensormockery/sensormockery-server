package stream

import (
	"log"
	"time"
)

// Stream represents a stream.
type Stream struct {
	ID         int
	WaveType   string
	Sensor     string
	NoiseCoeff float64
	BrokerURL  string
}

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
	for {
		if !currStreams[stream.ID] {
			delete(currStreams, stream.ID)
			return
		}

		log.Printf("%s\n", stream.WaveType)
		time.Sleep(time.Second)
	}
}
