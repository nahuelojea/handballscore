package coaches

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
	"github.com/nahuelojea/handballscore/storage"
)

func GetAvatar(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	coach, _, err := coaches_repository.GetCoach(Id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get coach: " + err.Error()
		return response
	}

	var filename = coach.Avatar
	if len(filename) < 1 {
		response.Status = http.StatusNotFound
		response.Message = "The coach has no avatar"
		return response
	}

	file, err := storage.GetFile(ctx, filename)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to download file in S3 " + err.Error()
		return response
	}

	response.CustomResp = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       file.String(),
		Headers: map[string]string{
			"Content-Type":        "application/octet-stream",
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", filename),
		},
	}
	return response
}
