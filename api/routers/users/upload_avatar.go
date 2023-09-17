package users

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/users_service"
	"github.com/nahuelojea/handballscore/storage"
)

func UploadAvatar(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest
	userId := claim.Id.Hex()

	var filename string
	var user models.User

	filename = "avatars/users/" + userId + ".jpg"
	user.Avatar = filename

	err := storage.UploadImage(ctx, request.Headers["Content-Type"], request.Body, filename)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to upload image: " + err.Error()
		return response
	}

	status, err := users_service.UpdateUser(user, userId)
	if err != nil || !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update user: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
