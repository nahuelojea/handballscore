package users

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
)

func GetUser(request events.APIGatewayProxyRequest) models.RespApi {
	var r models.RespApi
	r.Status = 400

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		r.Message = "'id' param is mandatory"
		return r
	}

	user, err := repositories.GetUser(Id)
	if err != nil {
		r.Message = "Error to get user " + err.Error()
		return r
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		r.Status = 500
		r.Message = "Error formating user to JSON " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(jsonResponse)
	return r
}
