package referees

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/referees_service"
)

func DeleteReferee(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	_, err := referees_service.DeleteReferee(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to delete referee: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Referee deleted"
	return response
}
