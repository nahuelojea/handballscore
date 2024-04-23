package matches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	matches_dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func SuspendMatch(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var suspendRequest matches_dto.SuspendRequest
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &suspendRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	_, err = matches_service.SuspendMatch(id, suspendRequest.Comments)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to suspend match: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match suspended"
	return response
}
