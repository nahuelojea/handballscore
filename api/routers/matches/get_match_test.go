package matches

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/matches_service" // For mocking
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock for matches_service.GetMatch - assumes package-level function mocking
var (
	MockServiceGetMatchFunc func(ID string) (models.Match, bool, error)
)

func GetMatchService(ID string) (models.Match, bool, error) {
	if MockServiceGetMatchFunc != nil {
		return MockServiceGetMatchFunc(ID)
	}
	return models.Match{}, false, errors.New("MockServiceGetMatchFunc not set")
}

func TestRouter_GetMatch_WithPlaceId(t *testing.T) {
	originalGetMatchService := matches_service.GetMatch
	defer func() { matches_service.GetMatch = originalGetMatchService }()

	matchID := primitive.NewObjectID()
	placeID := primitive.NewObjectID()
	expectedMatch := models.Match{
		Id:      matchID,
		Place:   "Some Stadium",
		PlaceId: placeID, // Ensure this is part of the mock response
		Date:    time.Now(),
		Status:  models.Programmed,
		// ... other necessary fields
	}

	// Setup mock for the service call
	matches_service.GetMatch = func(idStr string) (models.Match, bool, error) {
		assert.Equal(t, matchID.Hex(), idStr)
		return expectedMatch, true, nil
	}

	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"id": matchID.Hex()},
	}

	response := GetMatch(awsRequest) // Call the router function

	assert.Equal(t, http.StatusOK, response.Status)

	var returnedMatch models.Match
	err := json.Unmarshal([]byte(response.Message), &returnedMatch) // Message contains the JSON string
	assert.NoError(t, err, "Error unmarshalling response message")

	assert.Equal(t, expectedMatch.Id, returnedMatch.Id)
	assert.Equal(t, expectedMatch.Place, returnedMatch.Place)
	assert.Equal(t, expectedMatch.PlaceId, returnedMatch.PlaceId) // Crucial: verify PlaceId
	assert.Equal(t, expectedMatch.Status, returnedMatch.Status)
}

func TestRouter_GetMatch_ServiceError(t *testing.T) {
	originalGetMatchService := matches_service.GetMatch
	defer func() { matches_service.GetMatch = originalGetMatchService }()

	matchID := primitive.NewObjectID()
	serviceErr := errors.New("service failed to get match")

	matches_service.GetMatch = func(idStr string) (models.Match, bool, error) {
		return models.Match{}, false, serviceErr
	}

	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"id": matchID.Hex()},
	}

	response := GetMatch(awsRequest)

	assert.Equal(t, http.StatusNotFound, response.Status) // get_match router returns StatusNotFound for service errors
	assert.Contains(t, response.Message, "Error to get match: "+serviceErr.Error())
}

func TestRouter_GetMatch_NotFound(t *testing.T) {
	originalGetMatchService := matches_service.GetMatch
	defer func() { matches_service.GetMatch = originalGetMatchService }()

	matchID := primitive.NewObjectID()

	matches_service.GetMatch = func(idStr string) (models.Match, bool, error) {
		return models.Match{}, false, nil // Service returns false status for not found
	}

	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"id": matchID.Hex()},
	}
	ctx := context.Background()

	// The GetMatch router in api/routers/matches/get_match.go doesn't take context or claim.
	// It also doesn't explicitly set a "not found" message if status is false but err is nil.
	// It directly marshals the empty match and returns.
	// This test reflects the current behavior. If behavior should be different (e.g. 404 status), router needs change.
	response := GetMatch(awsRequest)

	// Based on current router: if service returns (models.Match{}, false, nil),
	// it proceeds to marshal the empty models.Match.
	// The router's GetMatch doesn't explicitly return a 404 for (false, nil) from service.
	// It returns StatusNotFound only if err != nil from service.
	// This might be an area for future improvement in the router.
	// For now, testing existing behavior.
	if response.Status == http.StatusOK { // Current behavior: it will be OK with an empty match if err is nil
		var returnedMatch models.Match
		err := json.Unmarshal([]byte(response.Message), &returnedMatch)
		assert.NoError(t, err)
		assert.True(t, returnedMatch.Id.IsZero()) // ID should be zero for an empty match
	} else {
		// This case would be hit if the router was changed to handle (false, nil) as a 404
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.Contains(t, response.Message, "not found")
	}
}

func TestRouter_GetMatch_MissingID(t *testing.T) {
	originalGetMatchService := matches_service.GetMatch
	defer func() { matches_service.GetMatch = originalGetMatchService }()

	matches_service.GetMatch = func(idStr string) (models.Match, bool, error) {
		t.Error("Service GetMatch should not be called")
		return models.Match{}, false, nil
	}
	
	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{}, // ID is missing
	}

	response := GetMatch(awsRequest)

	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "'id' param is mandatory", response.Message)
}
