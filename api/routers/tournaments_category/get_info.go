package tournaments

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/tournaments_info_service"
)

func GetInfo(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]

	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	tournamentInfo, err := tournaments_info_service.GetInfo(id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get tournament info: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(tournamentInfo)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating tournament info to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
