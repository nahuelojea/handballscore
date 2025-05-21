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
	matches_dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/services/matches_service" // For mocking
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMatchesService is a mock type for the MatchesService type
type MockMatchesService struct {
	mock.Mock
}

// This is a simplified way to allow mocking package-level functions from matches_service.
// In a real scenario, matches_service would ideally have an interface.
var (
	MockServiceProgramMatchFunc func(programMatchRequest matches_dto.ProgramMatchRequest, id string) (bool, error)
)

// Override actual service call for testing
func ProgramMatchService(programMatchRequest matches_dto.ProgramMatchRequest, id string) (bool, error) {
	if MockServiceProgramMatchFunc != nil {
		return MockServiceProgramMatchFunc(programMatchRequest, id)
	}
	return false, errors.New("MockServiceProgramMatchFunc not set")
}

func TestRouter_ProgramMatch_WithPlaceId(t *testing.T) {
	// Store original and defer restore
	originalProgramMatchService := matches_service.ProgramMatch
	defer func() { matches_service.ProgramMatch = originalProgramMatchService }()

	matchID := primitive.NewObjectID()
	placeID := primitive.NewObjectID()

	programRequest := matches_dto.ProgramMatchRequest{
		Date:    time.Now().Add(72 * time.Hour),
		Place:   "Grand Stadium",
		PlaceId: placeID, // Crucial: test this is passed
	}

	// Setup mock for the service call
	matches_service.ProgramMatch = func(req matches_dto.ProgramMatchRequest, idStr string) (bool, error) {
		assert.Equal(t, matchID.Hex(), idStr)
		assert.Equal(t, programRequest.Date.Unix(), req.Date.Unix())
		assert.Equal(t, programRequest.Place, req.Place)
		assert.Equal(t, programRequest.PlaceId, req.PlaceId) // Verify PlaceId in service call
		return true, nil
	}

	bodyBytes, _ := json.Marshal(programRequest)
	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"id": matchID.Hex()},
		Body:                  string(bodyBytes), // Body is passed via context in main handler
	}

	// Simulate context value for body as done in main_handler.go
	ctx := context.WithValue(context.Background(), dto.Key("body"), string(bodyBytes))

	response := ProgramMatch(ctx, awsRequest) // Call the router function

	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "Match programmed", response.Message)
}

func TestRouter_ProgramMatch_ServiceError(t *testing.T) {
	originalProgramMatchService := matches_service.ProgramMatch
	defer func() { matches_service.ProgramMatch = originalProgramMatchService }()

	matchID := primitive.NewObjectID()
	serviceErr := errors.New("service failed to program")

	programRequest := matches_dto.ProgramMatchRequest{Date: time.Now()}

	matches_service.ProgramMatch = func(req matches_dto.ProgramMatchRequest, idStr string) (bool, error) {
		return false, serviceErr
	}

	bodyBytes, _ := json.Marshal(programRequest)
	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"id": matchID.Hex()},
	}
	ctx := context.WithValue(context.Background(), dto.Key("body"), string(bodyBytes))

	response := ProgramMatch(ctx, awsRequest)

	assert.Equal(t, http.StatusInternalServerError, response.Status)
	assert.Contains(t, response.Message, "Error to program match data: "+serviceErr.Error())
}

func TestRouter_ProgramMatch_MissingID(t *testing.T) {
	originalProgramMatchService := matches_service.ProgramMatch
	defer func() { matches_service.ProgramMatch = originalProgramMatchService }()

	matches_service.ProgramMatch = func(req matches_dto.ProgramMatchRequest, idStr string) (bool, error) {
		// Should not be called
		t.Error("Service ProgramMatch should not be called when ID is missing")
		return false, nil
	}
	
	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{}, // ID is missing
	}
	ctx := context.Background() // Body not relevant here

	response := ProgramMatch(ctx, awsRequest)

	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "'id' param is mandatory", response.Message)
}

func TestRouter_ProgramMatch_InvalidJSON(t *testing.T) {
	originalProgramMatchService := matches_service.ProgramMatch
	defer func() { matches_service.ProgramMatch = originalProgramMatchService }()

	matches_service.ProgramMatch = func(req matches_dto.ProgramMatchRequest, idStr string) (bool, error) {
		// Should not be called
		t.Error("Service ProgramMatch should not be called with invalid JSON")
		return false, nil
	}

	matchID := primitive.NewObjectID()
	awsRequest := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"id": matchID.Hex()},
	}
	// Simulate context value for body as done in main_handler.go
	ctx := context.WithValue(context.Background(), dto.Key("body"), "{invalid json")

	response := ProgramMatch(ctx, awsRequest)
	
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Contains(t, response.Message, "Invalid data format")
}
