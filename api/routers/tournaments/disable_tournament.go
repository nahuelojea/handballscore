package tournaments

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/tournaments_repository"
)

func DisableTournament(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	isDisabled, err := tournaments_repository.DisableTournament(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable tournament: " + err.Error()
		return response
	}

	if !isDisabled {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable tournament in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Tournament disabled"
	return response
}
