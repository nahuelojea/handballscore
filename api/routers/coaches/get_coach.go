package coaches

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/coaches_service"
)

func GetCoach(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	coach, _, err := coaches_service.GetCoach(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get coach: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(coach)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating coach to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
