package services

import (
	"errors"

	"github.com/TonyJ3/song-service/models"
	"github.com/TonyJ3/song-service/repository"
)

var songRepo repository.SongRepository

func SetRepository(repo repository.SongRepository) {
	songRepo = repo
}

func GetAllSongs() ([]models.Song, error) {
	if songRepo == nil {
		return nil, errors.New("repository is not initialized")
	}
	return songRepo.GetAllSongs() //repository.GetAllSongs()
}

func GetSongByID(id string) (models.Song, error) {
	return songRepo.GetSongByID(id) //repository.GetSongByID(id)
}

func CreateSong(song models.Song) (models.Song, error) {
	return songRepo.CreateSong(song) //repository.CreateSong(song)
}

func UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	return songRepo.UpdateSong(id, updatedSong) //repository.UpdateSong(id, updatedSong)
}

func DeleteSong(id string) error {
	return songRepo.DeleteSong(id) //repository.DeleteSong(id
}
