package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TonyJ3/song-service/api"
)

func main() {
	// Setup the router (mux)
	router := api.SetupRouter()

	// Notifaction that the system is running
	fmt.Println("Localhost:8080 is running")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
