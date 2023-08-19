package matches

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
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

	var programMatchRequest dto.ProgramMatchRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &programMatchRequest)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Invalid data format: " + err.Error()
		return response
	}

	if programMatchRequest.Date.Compare(time.Now()) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "The date cannot be earlier than the current date"
		return response
	}

	status, err := matches_repository.ProgramMatch(programMatchRequest.Date, programMatchRequest.Place, Id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to program match data: " + err.Error()
		return response
	}

	if !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to program match data"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Match programmed"
	return response
}
