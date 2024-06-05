package match_players

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	matches_dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/services/match_players_service"
)

func UpdateNumber(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var numberRequest matches_dto.NumberRequest
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &numberRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	if len(numberRequest.Number) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "Number is required"
		return response
	}

	_, err = match_players_service.UpdateNumber(id, numberRequest.Number)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update number: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Number updated"
	return response
}
