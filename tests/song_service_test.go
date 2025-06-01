package tests

import (
	"testing"

	"github.com/TonyJ3/song-service/models"
	"github.com/TonyJ3/song-service/repository"
	"github.com/TonyJ3/song-service/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateSong(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockSongRepository()
	services.SetRepository(mockRepo)

	inputSong := models.Song{
		ID:     primitive.NewObjectID(),
		Title:  "Test Song",
		Artist: "Test Artist",
	}
	expectedSong := inputSong

	mockRepo.On("CreateSong", inputSong).Return(expectedSong, nil)

	// Act
	result, err := services.CreateSong(inputSong)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedSong, result)
	mockRepo.AssertCalled(t, "CreateSong", inputSong)
	mockRepo.AssertExpectations(t)
}
