package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sensormockery/sensormockery-server/pkg/db/dto"
	streamUtil "github.com/sensormockery/sensormockery-server/pkg/stream"
)

const (
	// CreateStreamPath is the api v1 path for stream creation.
	CreateStreamPath = "stream"
	// DeleteStreamPath is the api v1 path for stream deletion.
	DeleteStreamPath = "deleteStream"
)

// Stream is a json representation of a stream.
type Stream struct {
	WaveType   string  `json:"wave_type"`
	Sensor     string  `json:"sensor"`
	NoiseCoeff float64 `json:"noise_coeff"`
	BrokerURL  string  `json:"broker_url"`
}

// CreateStreamResp is a json representation of a create stream response.
type CreateStreamResp struct {
	ID int `json:"id"`
}

func handleStreamCreation(w http.ResponseWriter, r *http.Request) {
	// Read request data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	var apiV1Stream Stream
	if err := json.Unmarshal(body, &apiV1Stream); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Write into DB
	streamDTO, err := dto.NewStream(apiV1Stream.WaveType, apiV1Stream.Sensor, apiV1Stream.BrokerURL, apiV1Stream.NoiseCoeff)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := streamDTO.Upload(); err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Start stream
	stream := &streamUtil.Stream{
		ID:         streamDTO.GetID(),
		WaveType:   apiV1Stream.WaveType,
		Sensor:     apiV1Stream.Sensor,
		NoiseCoeff: apiV1Stream.NoiseCoeff,
		BrokerURL:  apiV1Stream.BrokerURL,
	}

	streamUtil.StartStream(stream)
	log.Printf("Created stream-%d", streamDTO.GetID())

	// Write response
	resp := &CreateStreamResp{
		ID: streamDTO.GetID(),
	}
	respAsJSON, err := json.Marshal(resp)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, string(respAsJSON))
}

func handleStreamDeletion(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	var id *CreateStreamResp
	if err := json.Unmarshal(body, &id); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	streamUtil.StopStream(id.ID)
	log.Printf("Stopped stream-%d", id.ID)

	writeResponse(w, http.StatusOK, "")
}
