package teams

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/teams_service"
)

func AddTeam(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var team models.Team
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &team)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(team.Name) == 0 {
		restResponse.Message = "Name is required"
		return restResponse
	}

	id, _, err := teams_service.CreateTeam(claim.AssociationId, team)
	if err != nil {
		restResponse.Message = "Error to create team: " + err.Error()
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
