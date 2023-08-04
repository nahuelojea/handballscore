package teams

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
)

func UpdateTeam(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	var team models.Team

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &team)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
	}

	status, err := teams_repository.UpdateTeam(team, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update team: " + err.Error()
		return response
	}

	if !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update team in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Team updated"
	return response
}
