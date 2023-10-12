package coaches

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/coaches_service"
)

func UploadAvatar(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	normalizedHeaders := make(map[string]string)
	for key, value := range request.Headers {
		normalizedKey := strings.ToLower(key)
		normalizedHeaders[normalizedKey] = value
	}

	err := coaches_service.UploadAvatar(ctx, request.Headers["content-type"], request.Body, id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update coach avatar: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
