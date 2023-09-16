package referees

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
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
	var referee models.Referee

	filename = "avatars/referees/" + id + ".jpg"
	referee.Avatar = filename

	err := storage.UploadImage(ctx, request.Headers["Content-Type"], request.Body, filename)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to upload image: " + err.Error()
		return response
	}

	status, err := referees_repository.UpdateReferee(referee, id)
	if err != nil || !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update referee " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
