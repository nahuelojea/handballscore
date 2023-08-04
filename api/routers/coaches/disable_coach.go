package coaches

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
)

func DisableCoach(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	isDisabled, err := coaches_repository.DisableCoach(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable coach: " + err.Error()
		return response
	}

	if !isDisabled {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable coach in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Coach disabled"
	return response
}
