package tournaments

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/tournaments_service"
)

func UpdateTournament(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	var tournament models.Tournament

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &tournament)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
	}

	_, err = tournaments_service.UpdateTournament(tournament, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update tournament: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Tournament updated"
	return response
}

func DeleteTournament(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	_, err := tournaments_service.DeleteTournament(Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to delete tournament: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Tournament deleted"
	return response
}
