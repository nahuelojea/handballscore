package teams

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
)

func DisableTeam(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	isDisabled, err := teams_repository.DisableTeam(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable team: " + err.Error()
		return response
	}

	if !isDisabled {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable team in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Team disabled"
	return response
}
