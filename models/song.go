package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Song represents a song object
type Song struct {
	/*ID     string `json:"id"`*/
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `json:"title"`
	Artist string             `json:"artist"`
	Genre  string             `json:"genre"`
}
