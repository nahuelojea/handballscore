package matches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	MatchesDTO "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func StartMatch(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	var startMatchRequest MatchesDTO.StartMatchRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &startMatchRequest)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
		return response
	}

	if len(startMatchRequest.PlayersLocal) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "There must be a minimum of one player on the home team"
		return response
	}

	if len(startMatchRequest.PlayersVisiting) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "There must be a minimum of one player on the visiting team"
		return response
	}

	if len(startMatchRequest.Referees) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "There must be a minimum of one referee"
		return response
	}

	if len(startMatchRequest.Timekeeper) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "Timekeeper is required"
		return response
	}

	if len(startMatchRequest.Scorekeeper) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "Scorekeeper is required"
		return response
	}

	_, err = matches_service.StartMatch(startMatchRequest, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to start match: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match started"
	return response
}
