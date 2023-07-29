package users

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

func GetUser(request events.APIGatewayProxyRequest) dto.RestResponse {
	var restResponse dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		restResponse.Status = 400
		restResponse.Message = "'id' param is mandatory"
		return restResponse
	}

	user, err := users_repository.GetUser(Id)
	if err != nil {
		restResponse.Status = 404
		restResponse.Message = "Error to get user: " + err.Error()
		return restResponse
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		restResponse.Status = 500
		restResponse.Message = "Error formating user to JSON " + err.Error()
		return restResponse
	}

	restResponse.Status = 200
	restResponse.Message = string(jsonResponse)
	return restResponse
}
