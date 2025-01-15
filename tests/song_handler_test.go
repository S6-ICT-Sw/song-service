package tests

//"bytes"
/*"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/gorilla/mux"
"github.com/stretchr/testify/assert"

"github.com/TonyJ3/song-service/api/song"
"github.com/TonyJ3/song-service/models"*/

/*func TestGetSongs(t *testing.T) {
	//song.Songs["1"] = song.Song{ID: "1", Title: "Song A", Artist: "Artist A", Genre: "Pop"}
	//song.Songs["2"] = song.Song{ID: "2", Title: "Song B", Artist: "Artist B", Genre: "Rock"}

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/songs", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Initialize the router and register the routes
	router := mux.NewRouter()
	song.RegisterRoutes(router)

	// Serve the HTTP request and record the response
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var songs []models.Song
	err = json.Unmarshal(rr.Body.Bytes(), &songs)
	assert.Nil(t, err)
}*/
