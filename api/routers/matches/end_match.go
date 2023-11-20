package matches

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func EndMatch(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	_, err := matches_service.EndMatch(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to end match: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match ended"
	return response
}
