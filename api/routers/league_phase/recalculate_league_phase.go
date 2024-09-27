package league_phase

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/league_phases_service"
)

func RecalculateLeaguePhase(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	err := league_phases_service.RecalculateTeamsScores(id)
	if err != nil {
		response.Message = "Error to recalculate league phase: " + err.Error()
		return response
	}

	response.Status = http.StatusCreated
	response.Message = "League phase recalculated"
	return response
}
