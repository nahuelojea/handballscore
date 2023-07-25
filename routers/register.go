package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
)

func Register(ctx context.Context) models.RespApi {
	var userModel models.User
	var restApiModel models.RespApi
	restApiModel.Status = 400

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &userModel)
	if err != nil {
		restApiModel.Message = err.Error()
		fmt.Println(restApiModel.Message)
		return restApiModel
	}

	if len(userModel.Email) == 0 {
		restApiModel.Message = "Email is mandatory"
		fmt.Println(restApiModel.Message)
		return restApiModel
	}
	if len(userModel.Password) < 6 {
		restApiModel.Message = "You must specify a password of at least 6 characters"
		fmt.Println(restApiModel.Message)
		return restApiModel
	}

	_, exist, _ := repositories.FindUserByEmail(userModel.Email)
	if exist {
		restApiModel.Message = "There is already a registered user with this email"
		fmt.Println(restApiModel.Message)
		return restApiModel
	}

	_, status, err := repositories.CreateUser(userModel)
	if err != nil {
		restApiModel.Message = "Error to register user: " + err.Error()
		fmt.Println(restApiModel.Message)
		return restApiModel
	}

	if !status {
		restApiModel.Message = "Error to register user"
		fmt.Println(restApiModel.Message)
		return restApiModel
	}

	restApiModel.Status = 200
	restApiModel.Message = "User registered"
	fmt.Println(restApiModel.Message)
	return restApiModel
}
