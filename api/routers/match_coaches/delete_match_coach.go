package match_coaches

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/match_coaches_service"
)

func DeleteMatchCoach(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	_, err := match_coaches_service.DeleteMatchCoach(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to delete match coach: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match coach deleted"
	return response
}
