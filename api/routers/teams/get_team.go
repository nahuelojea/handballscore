package teams

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/teams_service"
)

func GetTeam(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	team, _, err := teams_service.GetTeam(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get team: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(team)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating team to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
