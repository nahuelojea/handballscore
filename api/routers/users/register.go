package users

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/associations_repository"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

func Register(ctx context.Context) dto.RestResponse {
	var user models.User
	var restResponse dto.RestResponse
	restResponse.Status = 400

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		restResponse.Message = err.Error()
		fmt.Println(restResponse.Message)
		return restResponse
	}

	if len(user.Email) == 0 {
		restResponse.Message = "Email is mandatory"
		fmt.Println(restResponse.Message)
		return restResponse
	}
	if len(user.Password) < 6 {
		restResponse.Message = "You must specify a password of at least 6 characters"
		fmt.Println(restResponse.Message)
		return restResponse
	}
	if len(user.AssociationId) == 0 {
		restResponse.Message = "Association id is mandatory"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	_, exist, _ := associations_repository.GetAssociation(user.AssociationId)
	if !exist {
		restResponse.Message = "No association found with this id"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	_, exist, _ = users_repository.FindUserByEmail(user.Email)
	if exist {
		restResponse.Message = "There is already a registered user with this email"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	_, status, err := users_repository.CreateUser(user)
	if err != nil {
		restResponse.Message = "Error to register user: " + err.Error()
		fmt.Println(restResponse.Message)
		return restResponse
	}

	if !status {
		restResponse.Message = "Error to register user"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	restResponse.Status = 201
	restResponse.Message = "User registered"
	fmt.Println(restResponse.Message)
	return restResponse
}
