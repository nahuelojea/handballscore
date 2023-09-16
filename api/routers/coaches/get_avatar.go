package coaches

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/coaches_service"
)

func GetAvatar(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	file, filename, err := coaches_service.GetAvatar(id, ctx)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get Avatar: " + err.Error()
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
