package users

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/users_service"
)

func UploadAvatar(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest
	userId := claim.Id.Hex()

	normalizedHeaders := make(map[string]string)
	for key, value := range request.Headers {
		normalizedKey := strings.ToLower(key)
		normalizedHeaders[normalizedKey] = value
	}

	err := users_service.UploadAvatar(ctx, request.Headers["content-type"], request.Body, userId)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update user avatar: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
