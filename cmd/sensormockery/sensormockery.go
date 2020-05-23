package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	// EnvPort is the name of the env var for port.
	EnvPort = "PORT"

	// DefaultPort to be listened on.
	DefaultPort = "8080"
)

func main() {
	http.HandleFunc("/mock", mockHTTP)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", getPort()), nil))
}

func getPort() string {
	port := os.Getenv(EnvPort)

	if len(port) == 0 {
		port = DefaultPort
	}

	log.Printf("Using PORT %s", port)

	return port
}

func mockHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Mocked service!\n")
	log.Print("Request taken!")
}
