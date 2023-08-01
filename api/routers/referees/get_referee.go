package referees

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
)

func GetReferee(request events.APIGatewayProxyRequest) dto.RestResponse {
	var restResponse dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		restResponse.Status = http.StatusBadRequest
		restResponse.Message = "'id' param is mandatory"
		return restResponse
	}

	referee, err := referees_repository.GetReferee(Id)
	if err != nil {
		restResponse.Status = http.StatusNotFound
		restResponse.Message = "Error to get referee: " + err.Error()
		return restResponse
	}

	jsonResponse, err := json.Marshal(referee)
	if err != nil {
		restResponse.Status = http.StatusInternalServerError
		restResponse.Message = "Error formating referee to JSON " + err.Error()
		return restResponse
	}

	restResponse.Status = http.StatusOK
	restResponse.Message = string(jsonResponse)
	return restResponse
}
