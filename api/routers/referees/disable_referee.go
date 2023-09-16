package referees

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/referees_service"
)

func DisableReferee(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	isDisabled, err := referees_service.DisableReferee(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable referee: " + err.Error()
		return response
	}

	if !isDisabled {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to disable referee in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Referee disabled"
	return response
}
