package match_coaches

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/match_coaches_service"
)

func GetMatchCoach(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	match, _, err := match_coaches_service.GetMatchCoach(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get match coach: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(match)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating match to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
