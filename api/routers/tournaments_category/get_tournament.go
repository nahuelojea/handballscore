package tournaments

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func GetTournamentCategory(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	tournament, _, err := tournaments_service.GetTournamentCategory(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get tournament category: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(tournament)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating tournament category to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
