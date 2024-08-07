package matches

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	MatchesDTO "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func AssignReferees(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	var assignRefereesRequest MatchesDTO.AssingRefereesRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &assignRefereesRequest)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
		return response
	}

	if len(assignRefereesRequest.Referees) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "Referees is required"
		return response
	}

	_, err = matches_service.AssingReferees(Id, assignRefereesRequest)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to assign referees: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Referees assigned"
	return response
}
