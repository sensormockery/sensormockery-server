package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sensormockery/sensormockery-server/pkg/db/dto"
)

const (
	// CreateStreamPath is the api v1 path for create stream.
	CreateStreamPath = "stream"
)

// Stream is a json representation of a stream.
type Stream struct {
	WaveType   string  `json:"wave_type"`
	Sensor     string  `json:"sensor"`
	NoiseCoeff float64 `json:"noise_coeff"`
}

// CreateStreamResp is a json representation of a create stream response.
type CreateStreamResp struct {
	ID int `json:"id"`
}

func handleStreamCreation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	var stream Stream
	if err := json.Unmarshal(body, &stream); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	streamDTO, err := dto.NewStream(stream.WaveType, stream.Sensor, stream.NoiseCoeff)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := streamDTO.Upload(); err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

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
