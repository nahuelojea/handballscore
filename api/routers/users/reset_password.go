package users

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/users_service"
)

func ResetPassword(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	err := users_service.ResetPassword(id)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to reset password: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Reset password successfully"
	return response
}
