package song

import (
	"encoding/json"
	"net/http"

	"github.com/TonyJ3/song-service/models"
	"github.com/TonyJ3/song-service/services"
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the routes for song resource
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/songs", GetSongs).Methods("GET")
	r.HandleFunc("/songs/{id}", GetSongByID).Methods("GET")
	r.HandleFunc("/songs", CreateSong).Methods("POST")
	r.HandleFunc("/songs/{id}", UpdateSong).Methods("PUT")
	r.HandleFunc("/songs/{id}", DeleteSong).Methods("DELETE")
}

func GetSongs(w http.ResponseWriter, r *http.Request) {
	songs := services.GetAllSongs()

	if err := json.NewEncoder(w).Encode(songs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetSongByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	song, err := services.GetSongByID(id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(song)
}

func CreateSong(w http.ResponseWriter, r *http.Request) {
	var newSong models.Song

	// Create a new JSON decoder and disallow unknown fields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newSong); err != nil {
		http.Error(w, "Invalid JSON format or extra fields", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newSong.Title == "" || newSong.Artist == "" || newSong.Genre == "" {
		http.Error(w, "Missing required fields: title, artist, or genre", http.StatusBadRequest)
		return
	}

	createdSong := services.CreateSong(newSong)
	json.NewEncoder(w).Encode(createdSong)
}

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedSong models.Song

	// Create a new JSON decoder and disallow unknown fields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the incoming JSON request
	if err := decoder.Decode(&updatedSong); err != nil {
		http.Error(w, "Invalid JSON format or extra fields", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if updatedSong.Title == "" || updatedSong.Artist == "" || updatedSong.Genre == "" {
		http.Error(w, "Missing required fields: title, artist, or genre", http.StatusBadRequest)
		return
	}

	//Update the song in the service layer
	song, err := services.UpdateSong(id, updatedSong)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(song)
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := services.DeleteSong(id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
