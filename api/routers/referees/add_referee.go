package referees

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/referees_service"
)

func AddReferee(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var referee models.Referee
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &referee)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(referee.Name) == 0 {
		restResponse.Message = "Name is required"
		return restResponse
	}
	if len(referee.Surname) == 0 {
		restResponse.Message = "Surname is required"
		return restResponse
	}
	if len(referee.Dni) == 0 {
		restResponse.Message = "Dni is required"
		return restResponse
	}

	id, _, err := referees_service.CreateReferee(claim.AssociationId, referee)
	if err != nil {
		restResponse.Message = "Error to create referee: " + err.Error()
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
