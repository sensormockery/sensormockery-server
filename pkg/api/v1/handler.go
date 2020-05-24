package v1

import (
	"fmt"
	"net/http"
	"strings"
)

// APIPrefix is the url path of api v1.
const APIPrefix = "/api/v1/"

// Endpoint represents a single endpoint of api v1.
type Endpoint struct {
	Handler http.HandlerFunc
	Method  string
}

// EndpointRegistry is a collection of all the endpoints.
// Initialised on package import.
type EndpointRegistry map[string]*Endpoint

var endpoints EndpointRegistry

func init() {
	endpoints = map[string]*Endpoint{
		"stream": &Endpoint{
			Handler: handleStreamCreation,
			Method:  http.MethodPost,
		},
	}
}

// Handler func for api v1.
func Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, APIPrefix)
	endpoint, ok := endpoints[path]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Path not found."))

		return
	}

	if endpoint.Method != r.Method {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Expected a %s request but got a %s one\n", endpoint.Method, r.Method)
		w.Write([]byte(msg))

		return
	}

	endpoint.Handler(w, r)
}
