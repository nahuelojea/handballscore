package match_coaches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	matchesDTO "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/match_coaches_service"
)

func AddMatchCoach(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var matchCoachRequest matchesDTO.MatchCoachRequest
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &matchCoachRequest)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if claim.Role != models.AdminRole && claim.TeamId != matchCoachRequest.Team.Id {
		restResponse.Message = "You are not allowed to create a match coach in this team"
		return restResponse
	}

	if len(matchCoachRequest.MatchId) == 0 {
		restResponse.Message = "Match id is required"
		return restResponse
	}
	if len(matchCoachRequest.CoachId) == 0 {
		restResponse.Message = "Player id is required"
		return restResponse
	}
	if len(matchCoachRequest.Team.Id) == 0 {
		restResponse.Message = "Team id is required"
		return restResponse
	}

	id, status, err := match_coaches_service.CreateMatchCoach(claim.AssociationId, matchCoachRequest)
	if err != nil {
		restResponse.Message = "Error to create match player: " + err.Error()
		return restResponse
	}

	if !status {
		restResponse.Message = "Error to create match player"
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
