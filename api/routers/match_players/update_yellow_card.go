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

func UpdateYellowCard(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var yellowCardRequest matches_dto.YellowCardRequest
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &yellowCardRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	_, err = match_players_service.UpdateYellowCard(id, yellowCardRequest.Add)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update yellow card: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Yellow card updated"
	return response
}
