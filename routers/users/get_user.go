package users

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

func GetUser(request events.APIGatewayProxyRequest) models.RespApi {
	var response models.RespApi

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		response.Status = 400
		response.Message = "'id' param is mandatory"
		return response
	}

	user, err := users_repository.GetUser(Id)
	if err != nil {
		response.Status = 404
		response.Message = "Error to get user: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		response.Status = 500
		response.Message = "Error formating user to JSON " + err.Error()
		return response
	}

	response.Status = 200
	response.Message = string(jsonResponse)
	return response
}
