package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/users_service"
)

func UpdateUser(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	var user models.User

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
	}

	status, err := users_service.UpdateUser(user, claim.Id.Hex())
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update user: " + err.Error()
		return response
	}

	if !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update user in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "User updated"
	return response
}
