package match_players

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	matchesDTO "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/match_players_service"
)

func AddMatchPlayer(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var matchPlayerRequest matchesDTO.MatchPlayerRequest
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &matchPlayerRequest)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if claim.Role != models.AdminRole && claim.TeamId != matchPlayerRequest.Team.Id {
		restResponse.Message = "You are not allowed to create a match player in this team"
		return restResponse
	}

	if len(matchPlayerRequest.MatchId) == 0 {
		restResponse.Message = "Match id is required"
		return restResponse
	}
	if len(matchPlayerRequest.PlayerId) == 0 {
		restResponse.Message = "Player id is required"
		return restResponse
	}
	if len(matchPlayerRequest.Team.Id) == 0 {
		restResponse.Message = "Team id is required"
		return restResponse
	}
	if len(matchPlayerRequest.Number) == 0 {
		restResponse.Message = "Number is required"
		return restResponse
	}

	id, status, err := match_players_service.CreateMatchPlayer(claim.AssociationId, matchPlayerRequest)
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
