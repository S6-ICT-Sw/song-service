package services

import (
	"github.com/TonyJ3/song-service/models"
	"github.com/TonyJ3/song-service/repository"
)

/*func GetAllSongs() []models.Song {
	return repository.GetAllSongs()
}*/

func GetAllSongs() ([]models.Song, error) {
	return repository.GetAllSongs()
}

func GetSongByID(id string) (models.Song, error) {
	return repository.GetSongByID(id)
}

/*func CreateSong(song models.Song) models.Song {
	return repository.CreateSong(song)
}*/

func CreateSong(song models.Song) (models.Song, error) {
	return repository.CreateSong(song)
}

func UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	return repository.UpdateSong(id, updatedSong)
}

func DeleteSong(id string) error {
	return repository.DeleteSong(id)
}
