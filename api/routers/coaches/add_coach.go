package coaches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
)

func AddCoach(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var coach models.Coach
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &coach)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(coach.Name) == 0 {
		restResponse.Message = "Name is required"
		return restResponse
	}
	if len(coach.Surname) == 0 {
		restResponse.Message = "Surname is required"
		return restResponse
	}
	if len(coach.Dni) == 0 {
		restResponse.Message = "Dni is required"
		return restResponse
	}
	if len(coach.TeamId) == 0 {
		restResponse.Message = "Team id is mandatory"
		return restResponse
	}

	_, exist, _ := teams_repository.GetTeam(coach.TeamId)
	if !exist {
		restResponse.Message = "No team found with this id"
		return restResponse
	}

	_, exist, _ = coaches_repository.GetCoachByDni(coach.Dni)
	if exist {
		restResponse.Message = "There is already a registered coach with this dni"
		return restResponse
	}

	id, status, err := coaches_repository.CreateCoach(claim.AssociationId, coach)
	if err != nil {
		restResponse.Message = "Error to create coach: " + err.Error()
		return restResponse
	}

	if !status {
		restResponse.Message = "Error to create coach"
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
