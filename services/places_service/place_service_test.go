package places_service

import (
	"context"
	"errors"
	"testing"

	"github.com/nahuelojea/handballscore/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockPlaceRepository is a mock type for the PlaceRepository type
type MockPlaceRepository struct {
	mock.Mock
}

func (m *MockPlaceRepository) CreatePlace(ctx context.Context, place *models.Place) (string, bool, error) {
	args := m.Called(ctx, place)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockPlaceRepository) GetPlace(ctx context.Context, id string) (models.Place, bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Place), args.Bool(1), args.Error(2)
}

func (m *MockPlaceRepository) GetPlaces(ctx context.Context, filter bson.M, page, pageSize int) ([]models.Place, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	return args.Get(0).([]models.Place), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlaceRepository) UpdatePlace(ctx context.Context, id string, place models.Place) (bool, error) {
	args := m.Called(ctx, id, place)
	return args.Bool(0), args.Error(1)
}

func (m *MockPlaceRepository) DeletePlace(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestPlaceService_CreatePlace(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo} // Inject mock

	place := &models.Place{Name: "New Place"}
	expectedID := primitive.NewObjectID().Hex()

	mockRepo.On("CreatePlace", context.Background(), place).Return(expectedID, true, nil).Once()

	id, status, err := service.CreatePlace(context.Background(), place)

	assert.NoError(t, err)
	assert.True(t, status)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_CreatePlace_Error(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	place := &models.Place{Name: "New Place"}
	expectedError := errors.New("repository error")

	mockRepo.On("CreatePlace", context.Background(), place).Return("", false, expectedError).Once()

	_, status, err := service.CreatePlace(context.Background(), place)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, status)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_GetPlace(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	placeID := primitive.NewObjectID().Hex()
	expectedPlace := models.Place{Id: primitive.NewObjectIDFromHexIgnoreError(placeID), Name: "Found Place"}

	mockRepo.On("GetPlace", context.Background(), placeID).Return(expectedPlace, true, nil).Once()

	place, status, err := service.GetPlace(context.Background(), placeID)

	assert.NoError(t, err)
	assert.True(t, status)
	assert.Equal(t, expectedPlace, place)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_GetPlace_NotFound(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	placeID := primitive.NewObjectID().Hex()
	expectedError := errors.New("not found") // Or whatever error the repo returns for not found

	mockRepo.On("GetPlace", context.Background(), placeID).Return(models.Place{}, false, expectedError).Once()

	_, status, err := service.GetPlace(context.Background(), placeID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, status)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_GetPlaces(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	filter := bson.M{"name": "Test"}
	page, pageSize := 1, 10
	expectedPlaces := []models.Place{{Name: "Test Place 1"}, {Name: "Test Place 2"}}
	expectedTotal := int64(2)

	mockRepo.On("GetPlaces", context.Background(), filter, page, pageSize).Return(expectedPlaces, expectedTotal, nil).Once()

	places, total, err := service.GetPlaces(context.Background(), filter, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlaces, places)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_UpdatePlace(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	placeID := primitive.NewObjectID().Hex()
	placeToUpdate := models.Place{Name: "Updated Place"}

	mockRepo.On("UpdatePlace", context.Background(), placeID, placeToUpdate).Return(true, nil).Once()

	status, err := service.UpdatePlace(context.Background(), placeID, placeToUpdate)

	assert.NoError(t, err)
	assert.True(t, status)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_UpdatePlace_Error(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	placeID := primitive.NewObjectID().Hex()
	placeToUpdate := models.Place{Name: "Updated Place"}
	expectedError := errors.New("update failed")

	mockRepo.On("UpdatePlace", context.Background(), placeID, placeToUpdate).Return(false, expectedError).Once()

	status, err := service.UpdatePlace(context.Background(), placeID, placeToUpdate)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, status)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_DeletePlace(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	placeID := primitive.NewObjectID().Hex()

	mockRepo.On("DeletePlace", context.Background(), placeID).Return(true, nil).Once()

	status, err := service.DeletePlace(context.Background(), placeID)

	assert.NoError(t, err)
	assert.True(t, status)
	mockRepo.AssertExpectations(t)
}

func TestPlaceService_DeletePlace_Error(t *testing.T) {
	mockRepo := new(MockPlaceRepository)
	service := PlaceService{Repository: mockRepo}

	placeID := primitive.NewObjectID().Hex()
	expectedError := errors.New("delete failed")

	mockRepo.On("DeletePlace", context.Background(), placeID).Return(false, expectedError).Once()

	status, err := service.DeletePlace(context.Background(), placeID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, status)
	mockRepo.AssertExpectations(t)
}
