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
	json.NewEncoder(w).Encode(songs)
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

	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdSong := services.CreateSong(newSong)
	json.NewEncoder(w).Encode(createdSong)
}

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedSong models.Song
	_ = json.NewDecoder(r.Body).Decode(&updatedSong)
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
