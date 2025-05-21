package matches_service

import (
	"context"
	"errors"
	"testing"
	"time"

	matches_dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/matches_repository" // To mock its functions
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMatchesRepository is a mock for matches_repository functions
// Since repository functions are package-level, we need a way to swap them out.
// This can be done by defining function variables in the service or by using an interface.
// For this example, we'll assume we can set mock functions for what the service calls.
// A better approach would be for the service to use an interface for the repository.

type MockMatchesRepository struct {
	mock.Mock
}

// This is a simplified way to allow mocking package-level functions.
// In a real scenario, you'd use interfaces or dependency injection for the repository.
var (
	MockGetMatchFunc           func(ID string) (models.Match, bool, error)
	MockProgramMatchFunc       func(timeVal time.Time, place, streamingUrl, id string, placeId primitive.ObjectID) (bool, error)
	MockGetMatchHeaderViewFunc func(ID string) (models.MatchHeaderView, bool, error)
	MockGetMatchHeadersFunc    func(filterOptions matches_repository.GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error)
)

// Override actual repository calls with mocks for testing
func GetMatch(ID string) (models.Match, bool, error) {
	if MockGetMatchFunc != nil {
		return MockGetMatchFunc(ID)
	}
	return models.Match{}, false, errors.New("MockGetMatchFunc not set")
}

func ProgramMatchRepo(timeVal time.Time, place, streamingUrl, id string, placeId primitive.ObjectID) (bool, error) {
	if MockProgramMatchFunc != nil {
		return MockProgramMatchFunc(timeVal, place, streamingUrl, id, placeId)
	}
	return false, errors.New("MockProgramMatchFunc not set")
}

func GetMatchHeaderView(ID string) (models.MatchHeaderView, bool, error) {
	if MockGetMatchHeaderViewFunc != nil {
		return MockGetMatchHeaderViewFunc(ID)
	}
	return models.MatchHeaderView{}, false, errors.New("MockGetMatchHeaderViewFunc not set")
}

func GetMatchHeadersRepo(filterOptions matches_repository.GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
	if MockGetMatchHeadersFunc != nil {
		return MockGetMatchHeadersFunc(filterOptions)
	}
	return nil, 0, 0, errors.New("MockGetMatchHeadersFunc not set")
}

// Actual tests
func TestMatchService_ProgramMatch_WithPlaceId(t *testing.T) {
	originalGetMatch := matches_repository.GetMatch
	originalProgramMatch := matches_repository.ProgramMatch
	defer func() {
		matches_repository.GetMatch = originalGetMatch
		matches_repository.ProgramMatch = originalProgramMatch
	}()

	matchID := primitive.NewObjectID()
	placeID := primitive.NewObjectID()
	requestDTO := matches_dto.ProgramMatchRequest{
		Date:    time.Now().Add(48 * time.Hour),
		Place:   "Some Arena",
		PlaceId: placeID,
	}

	// Mock GetMatch (called before ProgramMatch by the service)
	matches_repository.GetMatch = func(ID string) (models.Match, bool, error) {
		assert.Equal(t, matchID.Hex(), ID)
		return models.Match{Id: matchID, AssociationId: "assoc1"}, true, nil
	}

	// Mock ProgramMatch (from repository)
	matches_repository.ProgramMatch = func(timeVal time.Time, place, streamingUrl, id string, pId primitive.ObjectID) (bool, error) {
		assert.Equal(t, matchID.Hex(), id)
		assert.Equal(t, requestDTO.Place, place)
		assert.Equal(t, requestDTO.PlaceId, pId) // Verify PlaceId is passed
		assert.Equal(t, requestDTO.Date.Unix(), timeVal.Unix())
		return true, nil
	}

	status, err := ProgramMatch(requestDTO, matchID.Hex())

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestMatchService_GetMatchesByJourney_PopulatesPlaceId(t *testing.T) {
	originalGetMatchHeaders := matches_repository.GetMatchHeaders
	defer func() { matches_repository.GetMatchHeaders = originalGetMatchHeaders }()

	assocID := "assocTest1"
	placeID := primitive.NewObjectID()
	matchHeaderView := models.MatchHeaderView{
		Id:        primitive.NewObjectID(),
		Place:     "Arena XYZ",
		PlaceId:   placeID,
		TeamHomeName: "Home",
		TeamAwayName: "Away",
		// ... other fields
	}

	matches_repository.GetMatchHeaders = func(filterOptions matches_repository.GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
		assert.Equal(t, assocID, filterOptions.AssociationId)
		return []models.MatchHeaderView{matchHeaderView}, 1, 1, nil
	}

	filterOpts := GetMatchesOptions{
		AssociationId: assocID,
		PageSize:      10, // Required for totalPages calculation
	}
	matchesJourney, total, totalPages, err := GetMatchesByJourney(filterOpts)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, totalPages)
	assert.Len(t, matchesJourney, 1)
	assert.Equal(t, matchHeaderView.Place, matchesJourney[0].Place)
	assert.Equal(t, placeID.Hex(), matchesJourney[0].PlaceId) // Verify PlaceId in DTO
}

func TestMatchService_GetMatchesByTeam_PopulatesPlaceId(t *testing.T) {
	originalGetMatchHeaders := matches_repository.GetMatchHeaders
	defer func() { matches_repository.GetMatchHeaders = originalGetMatchHeaders }()

	assocID := "assocTest2"
	teamID := primitive.NewObjectID()
	placeID := primitive.NewObjectID()

	matchHeaderView := models.MatchHeaderView{
		Id:           primitive.NewObjectID(),
		Place:        "Stadium ABC",
		PlaceId:      placeID,
		TeamHomeId:   teamID,
		TeamHomeName: "Team Alpha",
		TeamAwayName: "Team Beta",
		// ... other fields
	}

	matches_repository.GetMatchHeaders = func(filterOptions matches_repository.GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
		assert.Equal(t, assocID, filterOptions.AssociationId)
		assert.Equal(t, teamID.Hex(), filterOptions.Teams[0].TeamId)
		return []models.MatchHeaderView{matchHeaderView}, 1, 1, nil
	}

	filterOpts := GetMatchesOptions{
		AssociationId: assocID,
		Teams:         []models.TournamentTeamId{{TeamId: teamID.Hex()}},
		PageSize:      10, // Required for totalPages calculation
	}
	teamMatches, total, totalPages, err := GetMatchesByTeam(filterOpts)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, totalPages)
	assert.Len(t, teamMatches, 1)
	assert.Equal(t, matchHeaderView.Place, teamMatches[0].Place)
	assert.Equal(t, placeID.Hex(), teamMatches[0].PlaceId) // Verify PlaceId in DTO
}
