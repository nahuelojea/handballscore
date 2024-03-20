package matches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	matches_dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func UpdateTimeouts(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var timeoutRequest matches_dto.TimeoutRequest
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &timeoutRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	var tournamentTeamId models.TournamentTeamId = models.TournamentTeamId{
		TeamId:  timeoutRequest.TournamentTeamId.TeamId,
		Variant: timeoutRequest.TournamentTeamId.Variant,
	}

	_, err = matches_service.UpdateTimeouts(id, tournamentTeamId, timeoutRequest.Add, timeoutRequest.Time)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update timeout: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Timeout updated"
	return response
}
