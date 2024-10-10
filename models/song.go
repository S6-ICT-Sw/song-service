package models

// Song represents a song object
type Song struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Genre  string `json:"genre"`
}
