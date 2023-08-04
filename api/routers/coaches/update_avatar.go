package coaches

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
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
	var coach models.Coach

	filename = "avatars/coaches/" + id + ".jpg"
	coach.Avatar = filename

	hasError, response := storage.UploadImage(ctx, request, response, filename)
	if hasError {
		return response
	}

	status, err := coaches_repository.UpdateCoach(coach, id)
	if err != nil || !status {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to update coach " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
