package dto

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sensormockery/sensormockery-server/pkg/db"
)

// Stream is a DTO mapping to a row in streams table.
type Stream struct {
	db         *sql.DB
	id         int
	locked     int
	brokerURL  string
	waveType   string
	sensor     string
	noiseCoeff float64
	startDate  string
}

const (
	// WaveTypeMaxLength represents max strlen of waveType.
	WaveTypeMaxLength = 16
	// SensorMaxLength represents max strlen of sensor.
	SensorMaxLength = 64
	// BrokerURLMaxLength represents max strlen of broker url.
	BrokerURLMaxLength = 256
	// StreamsTable is the name of the table in db.
	StreamsTable = "streams"
	// Unlocked stream is one that is not currently emitting mocks.
	Unlocked = 0
	// Locked stream is one that is currently emitting mocks.
	Locked = 1
)

// NewStream returns a new Stream.
func NewStream(waveType, sensor, brokerURL string, noiseCoeff float64) (*Stream, error) {
	s := &Stream{}

	if err := s.SetWaveType(waveType); err != nil {
		return nil, err
	}

	if err := s.SetSensor(sensor); err != nil {
		return nil, err
	}

	if err := s.SetNoiseCoeff(noiseCoeff); err != nil {
		return nil, err
	}

	if err := s.SetBrokerURL(brokerURL); err != nil {
		return nil, err
	}

	dbConn, err := db.GetDBConn()
	if err != nil {
		return nil, fmt.Errorf("Error obtaining db connection: %s", err.Error())
	}
	s.db = dbConn

	return s, nil
}

// Upload a stream DTO to the db.
// After upload id is set the the id of the row inserted.
func (s *Stream) Upload() error {
	stmt, err := s.db.Prepare("INSERT INTO " + StreamsTable +
		"(locked, broker_url, wave_type, sensor, noise_coeff, start_date)" +
		" VALUES ($1, $2, $3, $4, $5, $6)" +
		" RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	currTime := time.Now().Format("2006-01-02 15:04:05")
	stmt.QueryRow(Unlocked, s.brokerURL, s.waveType, s.sensor, s.noiseCoeff, currTime).Scan(&(s.id))

	return nil
}

// Delete a stream DTO from the db.
func (s *Stream) Delete() error {
	return nil
}

// SetWaveType is a setter for waveType.
func (s *Stream) SetWaveType(waveType string) error {
	if len(waveType) > WaveTypeMaxLength {
		return fmt.Errorf("WaveType should shorter than %d", WaveTypeMaxLength)
	}

	s.waveType = waveType
	return nil
}

// SetSensor is a setter for sensor.
func (s *Stream) SetSensor(sensor string) error {
	if len(sensor) > SensorMaxLength {
		return fmt.Errorf("WaveType should shorter than %d", SensorMaxLength)
	}

	s.sensor = sensor
	return nil
}

// SetNoiseCoeff is a setter for noiseCoeff.
func (s *Stream) SetNoiseCoeff(noiseCoeff float64) error {
	if noiseCoeff < 0 {
		return fmt.Errorf("NoiseCoeff should be > 0")
	}

	s.noiseCoeff = noiseCoeff
	return nil
}

// SetBrokerURL is a setter for brokerURL.
func (s *Stream) SetBrokerURL(brokerURL string) error {
	if len(brokerURL) > BrokerURLMaxLength {
		return fmt.Errorf("BrokerURL should be shorter than %d", BrokerURLMaxLength)
	}

	s.brokerURL = brokerURL
	return nil
}

// GetID returns the ID of the stream DTO.
func (s *Stream) GetID() int {
	return s.id
}
