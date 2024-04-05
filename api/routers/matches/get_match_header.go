package matches

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func GetMatchHeader(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	matchHeader, _, err := matches_service.GetMatchHeader(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get match header: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(matchHeader)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating match header to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
