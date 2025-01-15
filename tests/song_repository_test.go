package tests

//"context"
/*"testing"

"github.com/TonyJ3/song-service/models"
"github.com/TonyJ3/song-service/repository"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
"go.mongodb.org/mongo-driver/mongo"*/

/*func TestCreatSong(t *testing.T) {
	// Create a mock collection
	mockCollection := new(repository.MockSongRepository)

	// Setup the expected behavior
	mockSong := models.Song{
		Title:  "Test Title",
		Artist: "Test Artist",
		Genre:  "Test Genre",
	}

	mockResult := &mongo.InsertOneResult{}
	mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(mockResult, nil)

	// Initialize the repository with the mock
	repository.NewMongoSongRepository(mockCollection)

	// Call the repository method
	createdSong, err := repository.CreateSong(mockSong)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockSong.Title, createdSong.Title)
	assert.Equal(t, mockSong.Artist, createdSong.Artist)

	// Verify the mock expectations
	mockCollection.AssertCalled(t, "InsertOne", mock.Anything, mock.Anything)
}*/
