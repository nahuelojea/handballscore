package tournaments

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/tournaments_service"
)

func CreateTournament(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var tournament models.Tournament
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &tournament)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(tournament.Name) == 0 {
		restResponse.Message = "Name is required"
		return restResponse
	}

	id, _, err := tournaments_service.CreateTournament(claim.AssociationId, tournament)

	if err != nil {
		restResponse.Message = "Error to create tournament: " + err.Error()
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
