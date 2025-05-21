package places

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetPlaces_Success(t *testing.T) {
	mockService := new(MockPlaceService) // Assuming MockPlaceService is accessible
	claim := models.Claim{AssociationId: "assocTest123"}
	expectedPlaces := []models.Place{
		{Name: "Place 1", AssociationId: claim.AssociationId},
		{Name: "Place 2", AssociationId: claim.AssociationId},
	}
	expectedTotalRecords := int64(2)

	page := 1
	pageSize := 10
	filter := bson.M{"association_id": claim.AssociationId} // Base filter

	mockService.On("GetPlaces", mock.Anything, filter, page, pageSize).Return(expectedPlaces, expectedTotalRecords, nil).Once()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"page":     strconv.Itoa(page),
			"pageSize": strconv.Itoa(pageSize),
		},
	}

	var response models.Response
	response = GetPlaces(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Places list", response.Message)
	assert.Equal(t, expectedTotalRecords, response.Total)

	var data []models.Place
	err := json.Unmarshal([]byte(response.Data.(string)), &data)
	if err != nil { // If Data is not stringified JSON
		data = response.Data.([]models.Place)
	}
	assert.Len(t, data, len(expectedPlaces))
	assert.Equal(t, expectedPlaces[0].Name, data[0].Name)

	mockService.AssertExpectations(t)
}

func TestGetPlaces_WithNameFilter(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}
	expectedPlaces := []models.Place{{Name: "Filtered Place", AssociationId: claim.AssociationId}}
	expectedTotalRecords := int64(1)

	page := 1
	pageSize := 5
	nameFilter := "Filtered"
	// Construct expected filter for the mock
	expectedFilter := bson.M{
		"association_id": claim.AssociationId,
		"name":           bson.M{"$regex": nameFilter, "$options": "i"},
	}

	mockService.On("GetPlaces", mock.Anything, expectedFilter, page, pageSize).Return(expectedPlaces, expectedTotalRecords, nil).Once()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"page":     strconv.Itoa(page),
			"pageSize": strconv.Itoa(pageSize),
			"name":     nameFilter,
		},
	}

	var response models.Response
	response = GetPlaces(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Places list", response.Message)
	assert.Equal(t, expectedTotalRecords, response.Total)
	var data []models.Place
	err := json.Unmarshal([]byte(response.Data.(string)), &data)
	if err != nil {
		data = response.Data.([]models.Place)
	}
	assert.Len(t, data, 1)
	assert.Equal(t, "Filtered Place", data[0].Name)

	mockService.AssertExpectations(t)
}

func TestGetPlaces_DefaultPagination(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}
	expectedPlaces := []models.Place{}
	expectedTotalRecords := int64(0)

	// Expected default values
	defaultPage := 1
	defaultPageSize := 20
	filter := bson.M{"association_id": claim.AssociationId}

	mockService.On("GetPlaces", mock.Anything, filter, defaultPage, defaultPageSize).Return(expectedPlaces, expectedTotalRecords, nil).Once()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{}, // No pagination params
	}

	var response models.Response
	response = GetPlaces(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetPlaces_ServiceError(t *testing.T) {
	mockService := new(MockPlaceService)
	claim := models.Claim{AssociationId: "assocTest123"}
	serviceError := errors.New("db connection lost")

	page := 1
	pageSize := 10
	filter := bson.M{"association_id": claim.AssociationId}

	mockService.On("GetPlaces", mock.Anything, filter, page, pageSize).Return(nil, int64(0), serviceError).Once()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"page":     strconv.Itoa(page),
			"pageSize": strconv.Itoa(pageSize),
		},
	}
	var response models.Response
	response = GetPlaces(context.Background(), request, mockService, claim, response)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Message, "Error to get places: "+serviceError.Error())
	mockService.AssertExpectations(t)
}
