package tournaments

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func UpdateTournamentCategory(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	var tournament models.TournamentCategory

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

	status, err := tournaments_service.UpdateTournamentCategory(tournament, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update tournament category: " + err.Error()
		return response
	}

	if !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update tournament category in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Tournament category updated"
	return response
}
