package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Stream is a json representation of a stream.
type Stream struct {
	WaveType   string  `json:"wave_type"`
	Sensor     string  `json:"sensor"`
	NoiseCoeff float64 `json:"noise_coeff"`
}

func handleStreamCreation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var stream Stream
	if err := json.Unmarshal(body, &stream); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
