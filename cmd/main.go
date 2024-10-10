package main

import (
	"log"
	"net/http"

	"github.com/TonyJ3/song-service/api"
)

func main() {
	// Setup the router (mux)
	router := api.SetupRouter()

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
