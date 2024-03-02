package tournaments

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/nahuelojea/handballscore/dto"
	TournamentDTO "github.com/nahuelojea/handballscore/dto/tournaments"
	TournamentsService "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func CreateTournamentCategory(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var createTournamentCategoryRequest TournamentDTO.CreateTournamentCategoryRequest
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &createTournamentCategoryRequest)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(createTournamentCategoryRequest.CategoryId) == 0 {
		restResponse.Message = "Category id is required"
		return restResponse
	}
	if len(createTournamentCategoryRequest.Teams) == 0 {
		restResponse.Message = "Teams are required"
		return restResponse
	}
	if reflect.DeepEqual(createTournamentCategoryRequest.LeaguePhase, TournamentDTO.LeaguePhaseRequest{}) &&
		reflect.DeepEqual(createTournamentCategoryRequest.PlayoffPhase, TournamentDTO.PlayoffPhaseRequest{}) &&
		reflect.DeepEqual(createTournamentCategoryRequest.LeagueAndPlayoff, TournamentDTO.LeagueAndPlayoffRequest{}) {
		restResponse.Message = "Tournament format is required"
		return restResponse
	}

	id, _, err := TournamentsService.CreateTournamentCategory(ctx, claim.AssociationId, createTournamentCategoryRequest)
	if err != nil {
		restResponse.Message = "Error to create tournament category: " + err.Error()
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
