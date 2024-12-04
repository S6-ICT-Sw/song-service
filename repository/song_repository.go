package repository

import (
	"errors"
	"strconv"
	"sync"

	"github.com/TonyJ3/song-service/models"
	//"github.com/TonyJ3/song-service/db"
)

var (
	songs = make(map[string]models.Song)
	//This is for 1 goroutine can access it at a time due to go is able to handel muti-request
	mutex  = &sync.Mutex{}
	nextID = 1
)

func GetAllSongs() []models.Song {
	mutex.Lock()
	defer mutex.Unlock()

	var allSongs []models.Song
	for _, song := range songs {
		allSongs = append(allSongs, song)
	}
	return allSongs
}

func GetSongByID(id string) (models.Song, error) {
	mutex.Lock()
	defer mutex.Unlock()

	song, exists := songs[id]
	if !exists {
		return models.Song{}, errors.New("Song not found")
	}

	return song, nil
}

func CreateSong(song models.Song) models.Song {
	mutex.Lock()
	defer mutex.Unlock()

	song.ID = strconv.Itoa(nextID)
	nextID++

	songs[song.ID] = song
	return song
}

func UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	mutex.Lock()
	defer mutex.Unlock()

	_, exists := songs[id]
	if !exists {
		return models.Song{}, errors.New("Song not found")
	}
	updatedSong.ID = id
	songs[id] = updatedSong
	return updatedSong, nil
}

func DeleteSong(id string) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, extist := songs[id]
	if !extist {
		return errors.New("Song not found")
	}
	delete(songs, id)
	return nil
}
