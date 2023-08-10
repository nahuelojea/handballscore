package players

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/players_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
)

func AddPlayer(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var player models.Player
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &player)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(player.Name) == 0 {
		restResponse.Message = "Name is required"
		return restResponse
	}
	if len(player.Surname) == 0 {
		restResponse.Message = "Surname is required"
		return restResponse
	}
	if len(player.Gender) == 0 {
		restResponse.Message = "Gender is required"
		return restResponse
	}
	if len(player.Dni) == 0 {
		restResponse.Message = "Dni is required"
		return restResponse
	}
	if len(player.TeamId) == 0 {
		restResponse.Message = "Team id is mandatory"
		return restResponse
	}

	_, exist, _ := teams_repository.GetTeam(player.TeamId)
	if !exist {
		restResponse.Message = "No team found with this id"
		return restResponse
	}

	_, exist, _ = players_repository.GetPlayerByDni(player.Dni)
	if exist {
		restResponse.Message = "There is already a registered player with this dni"
		return restResponse
	}

	id, status, err := players_repository.CreatePlayer(claim.AssociationId, player)
	if err != nil {
		restResponse.Message = "Error to create player: " + err.Error()
		return restResponse
	}

	if !status {
		restResponse.Message = "Error to create player"
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
