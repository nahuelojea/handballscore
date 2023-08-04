package teams

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
	"github.com/nahuelojea/handballscore/storage"
)

func UpdateAvatar(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	var filename string
	var team models.Team

	filename = "avatars/teams/" + id + ".jpg"
	team.Avatar = filename

	hasError, response := storage.UploadImage(ctx, request, response, filename)
	if hasError {
		return response
	}

	status, err := teams_repository.UpdateTeam(team, id)
	if err != nil || !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update team " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
