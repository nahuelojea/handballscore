package referees

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
)

func GetReferee(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	referee, err := referees_repository.GetReferee(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get referee: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(referee)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating referee to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
