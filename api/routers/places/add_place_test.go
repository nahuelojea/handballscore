package places

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockPlaceService is a mock type for the PlaceService type
type MockPlaceService struct {
	mock.Mock
}

func (m *MockPlaceService) CreatePlace(ctx context.Context, place *models.Place) (string, bool, error) {
	args := m.Called(ctx, place)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockPlaceService) GetPlace(ctx context.Context, id string) (models.Place, bool, error) {
	args := m.Called(ctx, id)
	var place models.Place
	if args.Get(0) != nil {
		place = args.Get(0).(models.Place)
	}
	return place, args.Bool(1), args.Error(2)
}

func (m *MockPlaceService) GetPlaces(ctx context.Context, filter primitive.M, page, pageSize int) ([]models.Place, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	var places []models.Place
	if args.Get(0) != nil {
		places = args.Get(0).([]models.Place)
	}
	return places, args.Get(1).(int64), args.Error(2)
}

func (m *MockPlaceService) UpdatePlace(ctx context.Context, id string, place models.Place) (bool, error) {
	args := m.Called(ctx, id, place)
	return args.Bool(0), args.Error(1)
}

func (m *MockPlaceService) DeletePlace(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestAddPlace_Success(t *testing.T) {
	mockService := new(MockPlaceService)
	placeToCreate := models.Place{
		Name: "Test Place",
		Ubication: models.UbicationCoordinates{
			Latitude:  10.0,
			Longitude: 20.0,
		},
	}
	expectedID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}

	// Configure the mock
	// We need to use a matcher for the place argument because CreatedDate and ModifiedDate are set inside CreatePlace (or repo)
	// and AssociationId is set in the router before calling the service.
	mockService.On("CreatePlace", mock.Anything, mock.MatchedBy(func(p *models.Place) bool {
		return p.Name == placeToCreate.Name && p.AssociationId == claim.AssociationId
	})).Return(expectedID.Hex(), true, nil).Once()

	bodyBytes, _ := json.Marshal(placeToCreate)
	request := events.APIGatewayProxyRequest{
		Body: string(bodyBytes),
	}

	var response models.Response
	response = AddPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Contains(t, response.Message, "Place created")
	var data dto.CreateIdResponse
	err := json.Unmarshal([]byte(response.Data.(string)), &data) // Assuming Data is stringified JSON
	if err != nil { // If Data is not stringified JSON, but the struct itself
		data = response.Data.(dto.CreateIdResponse)
	}
	assert.Equal(t, expectedID.Hex(), data.Id)
	mockService.AssertExpectations(t)
}

func TestAddPlace_EmptyName(t *testing.T) {
	mockService := new(MockPlaceService) // Service won't be called
	placeToCreate := models.Place{Name: ""} // Empty name
	claim := models.Claim{AssociationId: "assocTest123"}

	bodyBytes, _ := json.Marshal(placeToCreate)
	request := events.APIGatewayProxyRequest{
		Body: string(bodyBytes),
	}

	var response models.Response
	response = AddPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "Name is required", response.Message)
	mockService.AssertNotCalled(t, "CreatePlace")
}

func TestAddPlace_ServiceError(t *testing.T) {
	mockService := new(MockPlaceService)
	placeToCreate := models.Place{Name: "Test Place"}
	claim := models.Claim{AssociationId: "assocTest123"}
	serviceError := errors.New("service unavailable")

	mockService.On("CreatePlace", mock.Anything, mock.MatchedBy(func(p *models.Place) bool {
		return p.Name == placeToCreate.Name
	})).Return("", false, serviceError).Once()

	bodyBytes, _ := json.Marshal(placeToCreate)
	request := events.APIGatewayProxyRequest{
		Body: string(bodyBytes),
	}

	var response models.Response
	response = AddPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode) // Or specific error code if service distinguishes
	assert.Contains(t, response.Message, "Error to create place: "+serviceError.Error())
	mockService.AssertExpectations(t)
}

func TestAddPlace_JsonParseError(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}

	request := events.APIGatewayProxyRequest{
		Body: "{invalid json", // Malformed JSON
	}

	var response models.Response
	response = AddPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Message, "Error to parse place")
	mockService.AssertNotCalled(t, "CreatePlace")
}
