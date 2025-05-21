package places

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeletePlace_Success(t *testing.T) {
	mockService := new(MockPlaceService) // Assuming MockPlaceService is accessible
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}
	existingPlace := models.Place{Id: placeID, AssociationId: claim.AssociationId, Name: "To Delete"}

	// Mock GetPlace for permission check
	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(existingPlace, true, nil).Once()
	// Mock DeletePlace
	mockService.On("DeletePlace", mock.Anything, placeID.Hex()).Return(true, nil).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = DeletePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Place deleted", response.Message)
	mockService.AssertExpectations(t)
}

func TestDeletePlace_Forbidden(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "myAssoc"}
	existingPlaceFromOtherAssoc := models.Place{Id: placeID, AssociationId: "otherAssoc", Name: "To Delete"}

	// Mock GetPlace for permission check - returns place from different association
	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(existingPlaceFromOtherAssoc, true, nil).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = DeletePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, "User does not have permission to delete this place", response.Message)
	mockService.AssertNotCalled(t, "DeletePlace") // Ensure DeletePlace is not called
	mockService.AssertExpectations(t)
}

func TestDeletePlace_GetPlaceError(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}
	getError := errors.New("get error")

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(models.Place{}, false, getError).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = DeletePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode) // Based on current router logic
	assert.Contains(t, response.Message, "Error to get place for validation: "+getError.Error())
	mockService.AssertNotCalled(t, "DeletePlace")
	mockService.AssertExpectations(t)
}

func TestDeletePlace_ServiceErrorOnDelete(t *testing.T) {
	mockService := new(MockPlaceService)
	placeID := primitive.NewObjectID()
	claim := models.Claim{AssociationId: "assocTest123"}
	existingPlace := models.Place{Id: placeID, AssociationId: claim.AssociationId, Name: "To Delete"}
	deleteError := errors.New("service delete failed")

	mockService.On("GetPlace", mock.Anything, placeID.Hex()).Return(existingPlace, true, nil).Once()
	mockService.On("DeletePlace", mock.Anything, placeID.Hex()).Return(false, deleteError).Once()

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": placeID.Hex()},
	}

	var response models.Response
	response = DeletePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode) // Or specific error code
	assert.Contains(t, response.Message, "Error to delete place: "+deleteError.Error())
	mockService.AssertExpectations(t)
}

func TestDeletePlace_MissingID(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{}, // No ID
	}

	var response models.Response
	response = DeletePlace(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "Id is required", response.Message)
	mockService.AssertNotCalled(t, "GetPlace")
	mockService.AssertNotCalled(t, "DeletePlace")
}
