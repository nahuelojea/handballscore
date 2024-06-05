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

func ProgramMatch(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	var programMatchRequest MatchesDTO.ProgramMatchRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &programMatchRequest)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
		return response
	}

	if len(programMatchRequest.AuthorizationCode) != 6 {
		response.Status = http.StatusBadRequest
		response.Message = "Authorization code must have 6 characters"
		return response
	}

	_, err = matches_service.ProgramMatch(programMatchRequest.Date, programMatchRequest.Place, programMatchRequest.AuthorizationCode, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to program match data: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match programmed"
	return response
}
