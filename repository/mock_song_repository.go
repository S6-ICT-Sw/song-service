package repository

import (
	//"context"

	"github.com/TonyJ3/song-service/models"

	"github.com/stretchr/testify/mock"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type MockSongRepository struct {
	mock.Mock
}

func NewMockSongRepository() *MockSongRepository {
	return &MockSongRepository{}
}

func (m *MockSongRepository) GetAllSongs() ([]models.Song, error) {
	args := m.Called()
	return args.Get(0).([]models.Song), args.Error(1)
}

/*func (m *MockSongRepository) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}*/

func (m *MockSongRepository) GetSongByID(id string) (models.Song, error) {
	args := m.Called(id)
	return args.Get(0).(models.Song), args.Error(1)
}

/*func (m *MockSongRepository) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}*/

func (m *MockSongRepository) CreateSong(song models.Song) (models.Song, error) {
	args := m.Called(song)
	return args.Get(0).(models.Song), args.Error(1)
}

/*func (m *MockSongRepository) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}*/

func (m *MockSongRepository) UpdateSong(id string, updatedSong models.Song) (models.Song, error) {
	args := m.Called(id, updatedSong)
	return args.Get(0).(models.Song), args.Error(1)
}

/*func (m *MockSongRepository) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, replacement)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}*/

func (m *MockSongRepository) DeleteSong(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

/*func (m *MockSongRepository) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}*/
