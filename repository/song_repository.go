package repository

import (
	"context"
	"errors"
	"time"

	"github.com/TonyJ3/song-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

/*var (
	songs = make(map[string]models.Song)
	//This is for 1 goroutine can access it at a time due to go is able to handel muti-request
	mutex  = &sync.Mutex{}
	nextID = 1
)*/

var songCollection *mongo.Collection

func InitRepository(client *mongo.Client, dbName, collectionName string) {
	songCollection = client.Database(dbName).Collection(collectionName)
}

/*func GetAllSongs() []models.Song {
	mutex.Lock()
	defer mutex.Unlock()

	var allSongs []models.Song
	for _, song := range songs {
		allSongs = append(allSongs, song)
	}
	return allSongs
}*/

func GetAllSongs() ([]models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := songCollection.Find(ctx, bson.D{})
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

/*func GetSongByID(id string) (models.Song, error) {
	mutex.Lock()
	defer mutex.Unlock()

	song, exists := songs[id]
	if !exists {
		return models.Song{}, errors.New("song not found")
	}

	return song, nil
}*/

func GetSongByID(id string) (models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Song{}, errors.New("invalid ID format")
	}

	var song models.Song
	if err := songCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&song); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Song{}, errors.New("song not found")
		}
		return models.Song{}, err
	}
	return song, nil
}

/*func CreateSong(song models.Song) models.Song {
	mutex.Lock()
	defer mutex.Unlock()

	song.ID = strconv.Itoa(nextID)
	nextID++

	songs[song.ID] = song
	return song
}*/

func CreateSong(song models.Song) (models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	song.ID = primitive.NewObjectID()
	_, err := songCollection.InsertOne(ctx, song)
	if err != nil {
		return models.Song{}, err
	}
	return song, nil
}

/*func UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	mutex.Lock()
	defer mutex.Unlock()

	_, exists := songs[id]
	if !exists {
		return models.Song{}, errors.New("song not found")
	}
	updatedSong.ID = id
	songs[id] = updatedSong
	return updatedSong, nil
}*/

func UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Song{}, errors.New("invalid ID format")
	}

	updatedSong.ID = objID
	_, err = songCollection.ReplaceOne(ctx, bson.M{"_id": objID}, updatedSong)
	if err != nil {
		return models.Song{}, err
	}
	return updatedSong, nil
}

/*func DeleteSong(id string) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, extist := songs[id]
	if !extist {
		return errors.New("song not found")
	}
	delete(songs, id)
	return nil
}*/

func DeleteSong(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	_, err = songCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
