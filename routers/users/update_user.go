package users

import (
	"context"
	"encoding/json"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

func UpdateUser(ctx context.Context, claim models.Claim) models.RespApi {
	var response models.RespApi

	var user models.User

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Status = 400
		response.Message = "Invalid data format: " + err.Error()
	}

	status, err := users_repository.UpdateUser(user, claim.Id.Hex())
	if err != nil {
		response.Status = 500
		response.Message = "Error to update user: " + err.Error()
		return response
	}

	if !status {
		response.Status = 500
		response.Message = "Error to update user in database"
		return response
	}

	response.Status = 200
	response.Message = "User updated"
	return response
}
