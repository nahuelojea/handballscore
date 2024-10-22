package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	userDTO "github.com/nahuelojea/handballscore/dto/users"
	"github.com/nahuelojea/handballscore/services/users_service"
)

func ChangePassword(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var changePasswordRequest userDTO.ChangePasswordRequest

	userId := claim.Id.Hex()
	
	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &changePasswordRequest)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
	}

	err = users_service.ChangePassword(userId, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to change password: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Password changed successfully"
	return response
}
