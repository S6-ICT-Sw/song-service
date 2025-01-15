package song

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/TonyJ3/song-service/models"
	"github.com/TonyJ3/song-service/services"
	"github.com/gorilla/mux"

	"github.com/TonyJ3/song-service/messaging"
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
	songs, err := services.GetAllSongs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(songs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetSongByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	song, err := services.GetSongByID(id) //services.GetSongByID(id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	//json.NewEncoder(w).Encode(song)

	if err := json.NewEncoder(w).Encode(song); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateSong(w http.ResponseWriter, r *http.Request) {
	//log.Println("Received request to create a song")

	var newSong models.Song

	// Create a new JSON decoder and disallow unknown fields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newSong); err != nil {
		//log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid JSON format or extra fields", http.StatusBadRequest)
		return
	}
	//log.Printf("Parsed song from request: %+v", newSong)

	// Validate required fields
	if newSong.Title == "" || newSong.Artist == "" || newSong.Genre == "" {
		http.Error(w, "Missing required fields: title, artist, or genre", http.StatusBadRequest)
		return
	}

	createdSong, err := services.CreateSong(newSong) //services.CreateSong(newSong)
	if err != nil {
		log.Printf("Failed to create song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//log.Printf("Successfully created song: %+v", createdSong)
	//json.NewEncoder(w).Encode(createdSong)

	// Notify the song suggestion service
	/*err = notifySongSuggestionService(createdSong)
	if err != nil {
		http.Error(w, "Failed to notify song suggestion service: "+err.Error(), http.StatusInternalServerError)
		return
	}*/

	// Publish the event to RabbitMQ
	// Initialize RabbitMQ connection and channel
	err = messaging.InitRabbitMQ()
	if err != nil {
		http.Error(w, "Failed to connect to RabbitMQ", http.StatusInternalServerError)
		return
	}
	defer messaging.CloseRabbitMQ() // Ensure the channel and connection are closed after the function finishes

	// Publish the "created" event to RabbitMQ
	err = messaging.PublishMessage(
		messaging.GetChannel(), // Use the channel from InitRabbitMQ
		"created",              // eventType: "created" for a new song
		createdSong.ID.Hex(),   // song_ID: the ID of the created song
		createdSong.Title,      // title: the song's title
		createdSong.Artist,     // artist: the song's artist
	)
	if err != nil {
		http.Error(w, "Failed to publish song creation event", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(createdSong); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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
	song, err := services.UpdateSong(id, updatedSong) //services.UpdateSong(id, updatedSong)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	//json.NewEncoder(w).Encode(song)

	if err := json.NewEncoder(w).Encode(song); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := services.DeleteSong(id) //services.DeleteSong(id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Publish the event to RabbitMQ
	// Initialize RabbitMQ connection and channel
	err = messaging.InitRabbitMQ()
	if err != nil {
		http.Error(w, "Failed to connect to RabbitMQ", http.StatusInternalServerError)
		return
	}
	defer messaging.CloseRabbitMQ() // Ensure the channel and connection are closed after the function finishes

	// Publish the delete event to RabbitMQ
	err = messaging.PublishMessage(messaging.GetChannel(), "deleted", id, "", "")
	if err != nil {
		log.Printf("Failed to publish delete event for song ID %s: %v", id, err)
		http.Error(w, "Failed to publish delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
