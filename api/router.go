package api

import (
	"github.com/TonyJ3/song-service/api/song"
	"github.com/gorilla/mux"
)

// SetupRouter initializes the routes
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Register song routes
	song.RegisterRoutes(router)

	return router
}
