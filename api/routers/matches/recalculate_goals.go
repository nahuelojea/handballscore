package matches

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func RecalculateGoals(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	_, err := matches_service.RecalculateMatchGoals(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get match: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match goals recalculated"
	return response
}
