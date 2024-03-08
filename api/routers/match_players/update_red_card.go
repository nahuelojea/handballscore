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

func UpdateRedCard(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var redCardRequest matches_dto.RedCardRequest
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &redCardRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	_, err = match_players_service.UpdateRedCard(id, redCardRequest.Add)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update red card: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Red card updated"
	return response
}
