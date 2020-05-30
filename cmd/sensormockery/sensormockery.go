package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	apiV1 "github.com/sensormockery/sensormockery-server/pkg/api/v1"
	"github.com/sensormockery/sensormockery-server/pkg/db"
	"github.com/sensormockery/sensormockery-server/pkg/env"
	"github.com/sensormockery/sensormockery-server/pkg/stream"
)

func main() {
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	go func() {
		stream.StartStreamListener()
	}()

	http.HandleFunc("/api/v1/", apiV1.Handler)

	port := env.GetPort()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
