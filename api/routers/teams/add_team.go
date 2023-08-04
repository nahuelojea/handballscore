package teams

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/associations_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
)

func AddTeam(ctx context.Context) dto.RestResponse {
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
	if len(team.AssociationId) == 0 {
		restResponse.Message = "Association id is mandatory"
		return restResponse
	}

	_, exist, _ := associations_repository.GetAssociation(team.AssociationId)
	if !exist {
		restResponse.Message = "No association found with this id"
		return restResponse
	}

	id, status, err := teams_repository.CreateTeam(team)
	if err != nil {
		restResponse.Message = "Error to create team: " + err.Error()
		return restResponse
	}

	if !status {
		restResponse.Message = "Error to create team"
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
