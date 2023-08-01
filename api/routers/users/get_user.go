package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

func GetUser(request events.APIGatewayProxyRequest) dto.RestResponse {
	var restResponse dto.RestResponse

	Id := request.QueryStringParameters["id"]
	if len(Id) < 1 {
		restResponse.Status = http.StatusBadRequest
		restResponse.Message = "'id' param is mandatory"
		return restResponse
	}

	user, err := users_repository.GetUser(Id)
	if err != nil {
		restResponse.Status = http.StatusNotFound
		restResponse.Message = "Error to get user: " + err.Error()
		return restResponse
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		restResponse.Status = http.StatusInternalServerError
		restResponse.Message = "Error formating user to JSON " + err.Error()
		return restResponse
	}

	fmt.Println("GET USER: ", string(jsonResponse))

	restResponse.Status = http.StatusOK
	restResponse.Message = string(jsonResponse)
	return restResponse
}
