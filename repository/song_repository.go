package repository

import (
	"context"
	"errors"

	//"log"
	"time"

	"github.com/TonyJ3/song-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SongRepository interface {
	CreateSong(song models.Song) (models.Song, error)
	GetAllSongs() ([]models.Song, error)
	GetSongByID(id string) (models.Song, error)
	UpdateSong(id string, updatedSong models.Song) (models.Song, error)
	DeleteSong(id string) error
}

type MongoSongRepository struct {
	Collection *mongo.Collection
}

func NewMongoSongRepository(collection *mongo.Collection) *MongoSongRepository {
	return &MongoSongRepository{Collection: collection}
}

func (r *MongoSongRepository) GetAllSongs() ([]models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//cursor, err := songCollection.Find(ctx, bson.D{})
	cursor, err := r.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []models.Song
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func (r *MongoSongRepository) GetSongByID(id string) (models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Song{}, errors.New("invalid ID format")
	}

	var song models.Song
	if err := r.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&song); err != nil { // songCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&song); err != nil
		if err == mongo.ErrNoDocuments {
			return models.Song{}, errors.New("song not found")
		}
		return models.Song{}, err
	}
	return song, nil
}

func (r *MongoSongRepository) CreateSong(song models.Song) (models.Song, error) {
	//log.Printf("Inserting song into MongoDB: %+v", song)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	song.ID = primitive.NewObjectID()
	_, err := r.Collection.InsertOne(ctx, song) //songCollection.InsertOne(ctx, song)
	if err != nil {
		//log.Printf("Failed to insert song into MongoDB: %v", err)
		return models.Song{}, err
	}
	//log.Printf("Successfully inserted song with ID: %s", song.ID)
	return song, nil
}

func (r *MongoSongRepository) UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Song{}, errors.New("invalid ID format")
	}

	updatedSong.ID = objID
	_, err = r.Collection.ReplaceOne(ctx, bson.M{"_id": objID}, updatedSong) //songCollection.ReplaceOne(ctx, bson.M{"_id": objID}, updatedSong)
	if err != nil {
		return models.Song{}, err
	}
	return updatedSong, nil
}

func (r *MongoSongRepository) DeleteSong(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": objID}) //songCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
