package users

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
	"github.com/nahuelojea/handballscore/storage"
)

func UpdateAvatar(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest
	userId := claim.Id.Hex()

	var filename string
	var user models.User

	filename = "avatars/users/" + userId + ".jpg"
	user.Avatar = filename

	hasError, response := storage.UploadImage(ctx, request, response, filename)
	if hasError {
		return response
	}

	status, err := users_repository.UpdateUser(user, userId)
	if err != nil || !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update user " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
