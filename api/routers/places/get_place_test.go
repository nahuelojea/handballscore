package places

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetPlace_Success(t *testing.T) {
	mockService := new(MockPlaceService) // Assuming MockPlaceService is defined in add_place_test.go or a shared file
	placeID := primitive.NewObjectID()
	expectedPlace := models.Place{
		Id:            placeID,
		Name:          "Test Place",
		AssociationId: "assocTest123",
	}
	claim := models.Claim{AssociationId: "assocTest123"}

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(expectedPlace, true, nil).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = GetPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	var data models.Place
	err := json.Unmarshal([]byte(response.Data.(string)), &data)
	if err != nil {
		data = response.Data.(models.Place)
	}
	assert.Equal(t, expectedPlace.Id, data.Id)
	assert.Equal(t, expectedPlace.Name, data.Name)
	mockService.AssertExpectations(t)
}

func TestGetPlace_NotFound(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(models.Place{}, false, nil).Once() // Service returns false for status

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = GetPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "Place not found", response.Message)
	mockService.AssertExpectations(t)
}

func TestGetPlace_ServiceError(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}
	serviceError := errors.New("service error")

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(models.Place{}, false, serviceError).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = GetPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode) // Or specific error code
	assert.Contains(t, response.Message, "Error to get place: "+serviceError.Error())
	mockService.AssertExpectations(t)
}

func TestGetPlace_Forbidden(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	placeFromOtherAssociation := models.Place{
		Id:            placeID,
		Name:          "Belongs to Other",
		AssociationId: "otherAssoc", // Different AssociationId
	}
	claim := models.Claim{AssociationId: "myAssoc"} // User's association

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(placeFromOtherAssociation, true, nil).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = GetPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, "User does not have permission to access this place", response.Message)
	mockService.AssertExpectations(t)
}

func TestGetPlace_MissingID(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{}, // No ID
	}

	var response models.Response
	response = GetPlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "Id is required", response.Message)
	mockService.AssertNotCalled(t, "GetPlace")
}
