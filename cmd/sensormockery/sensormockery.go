package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sensormockery/sensormockery-server/pkg/db"
	"github.com/sensormockery/sensormockery-server/pkg/env"
)

func main() {
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/mock", mockHTTP)

	port := env.GetPort()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func mockHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Mocked service!\n")
	log.Print("Request taken!")
}
