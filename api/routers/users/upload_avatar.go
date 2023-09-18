package users

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/users_service"
)

func UploadAvatar(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest
	userId := claim.Id.Hex()

	err := users_service.UploadAvatar(ctx, request.Headers["Content-Type"], request.Body, userId)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update user avatar: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
