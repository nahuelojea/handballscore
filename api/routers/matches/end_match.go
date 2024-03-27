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

func EndMatch(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var endMatchRequest matches_dto.EndMatchRequest
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &endMatchRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	_, err = matches_service.EndMatch(id, endMatchRequest.Comments)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to end match: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match ended"
	return response
}
