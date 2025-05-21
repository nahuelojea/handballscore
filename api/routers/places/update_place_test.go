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

func TestUpdatePlace_Success(t *testing.T) {
	mockService := new(MockPlaceService) // Assuming MockPlaceService is accessible
	placeID := primitive.NewObjectID()
	placeToUpdate := models.Place{
		Name: "Updated Place Name",
		Ubication: models.UbicationCoordinates{
			Latitude:  15.0,
			Longitude: 25.0,
		},
	}
	claim := models.Claim{AssociationId: "assocTest123"}
	existingPlace := models.Place{Id: placeID, AssociationId: claim.AssociationId, Name: "Old Name"}

	// Mock GetPlace for permission check
	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(existingPlace, true, nil).Once()
	// Mock UpdatePlace
	mockService.On("UpdatePlace", mock.Anything, placeID.Hex(), mock.MatchedBy(func(p models.Place) bool {
		return p.Name == placeToUpdate.Name && p.AssociationId == claim.AssociationId && p.Id == placeID
	})).Return(true, nil).Once()

	bodyBytes, _ := json.Marshal(placeToUpdate)
	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
		Body:           string(bodyBytes),
	}

	var response models.Response
	response = UpdatePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Place updated", response.Message)
	mockService.AssertExpectations(t)
}

func TestUpdatePlace_Forbidden(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	placeToUpdate := models.Place{Name: "Updated Name"}
	claim := models.Claim{AssociationId: "myAssoc"}
	existingPlaceFromOtherAssoc := models.Place{Id: placeID, AssociationId: "otherAssoc", Name: "Old Name"}

	// Mock GetPlace for permission check - returns place from different association
	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(existingPlaceFromOtherAssoc, true, nil).Once()

	bodyBytes, _ := json.Marshal(placeToUpdate)
	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
		Body:           string(bodyBytes),
	}

	var response models.Response
	response = UpdatePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, "User does not have permission to update this place", response.Message)
	mockService.AssertNotCalled(t, "UpdatePlace") // Ensure UpdatePlace is not called
	mockService.AssertExpectations(t)
}

func TestUpdatePlace_GetPlaceError(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	placeToUpdate := models.Place{Name: "Updated Name"}
	claim := models.Claim{AssociationId: "assocTest123"}
	getError := errors.New("get error")

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(models.Place{}, false, getError).Once()

	bodyBytes, _ := json.Marshal(placeToUpdate)
	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
		Body:           string(bodyBytes),
	}

	var response models.Response
	response = UpdatePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode) // Based on current router logic
	assert.Contains(t, response.Message, "Error to get place for validation: "+getError.Error())
	mockService.AssertNotCalled(t, "UpdatePlace")
	mockService.AssertExpectations(t)
}

func TestUpdatePlace_ServiceErrorOnUpdate(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	placeToUpdate := models.Place{Name: "Updated Place Name"}
	claim := models.Claim{AssociationId: "assocTest123"}
	existingPlace := models.Place{Id: placeID, AssociationId: claim.AssociationId, Name: "Old Name"}
	updateError := errors.New("service update failed")

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(existingPlace, true, nil).Once()
	mockService.On("UpdatePlace", mock.Anything, placeID.Hex(), mock.MatchedBy(func(p models.Place) bool {
		return p.Name == placeToUpdate.Name
	})).Return(false, updateError).Once()

	bodyBytes, _ := json.Marshal(placeToUpdate)
	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
		Body:           string(bodyBytes),
	}

	var response models.Response
	response = UpdatePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode) // Or specific error code
	assert.Contains(t, response.Message, "Error to update place: "+updateError.Error())
	mockService.AssertExpectations(t)
}

func TestUpdatePlace_MissingID(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}
	placeToUpdate := models.Place{Name: "Updated Place Name"}

	bodyBytes, _ := json.Marshal(placeToUpdate)
	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{}, // No ID
		Body:           string(bodyBytes),
	}

	var response models.Response
	response = UpdatePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "Id is required", response.Message)
	mockService.AssertNotCalled(t, "GetPlace")
	mockService.AssertNotCalled(t, "UpdatePlace")
}

func TestUpdatePlace_JsonParseError(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
		Body:           "{invalid json", // Malformed JSON
	}

	var response models.Response
	response = UpdatePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Message, "Error to parse place")
	mockService.AssertNotCalled(t, "GetPlace")
	mockService.AssertNotCalled(t, "UpdatePlace")
}
