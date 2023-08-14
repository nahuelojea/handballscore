package tournaments

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/tournaments_repository"
)

func GetTournament(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	tournament, _, err := tournaments_repository.GetTournament(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get tournament: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(tournament)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating tournament to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
