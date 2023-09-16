package coaches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/coaches_service"
)

func UpdateCoach(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	var coach models.Coach

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &coach)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
	}

	status, err := coaches_service.UpdateCoach(coach, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update coach: " + err.Error()
		return response
	}

	if !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update coach in database"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Coach updated"
	return response
}
