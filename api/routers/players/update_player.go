package players

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/players_repository"
)

func UpdatePlayer(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	var player models.Player

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &player)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
	}

	status, err := players_repository.UpdatePlayer(player, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update player: " + err.Error()
		return response
	}

	if !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update player in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Player updated"
	return response
}
