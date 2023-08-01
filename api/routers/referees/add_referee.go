package referees

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/associations_repository"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
)

func AddReferee(ctx context.Context) dto.RestResponse {
	var referee models.Referee
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &referee)
	if err != nil {
		restResponse.Message = err.Error()
		fmt.Println(restResponse.Message)
		return restResponse
	}

	if len(referee.Name) == 0 {
		restResponse.Message = "Name is required"
		fmt.Println(restResponse.Message)
		return restResponse
	}
	if len(referee.Surname) == 0 {
		restResponse.Message = "Surname is required"
		fmt.Println(restResponse.Message)
		return restResponse
	}
	if len(referee.AssociationId) == 0 {
		restResponse.Message = "Association id is mandatory"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	_, exist, _ := associations_repository.GetAssociation(referee.AssociationId)
	if !exist {
		restResponse.Message = "No association found with this id"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	_, exist, _ = referees_repository.FindRefereeByDni(referee.Dni)
	if exist {
		restResponse.Message = "There is already a registered referee with this dni"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	id, status, err := referees_repository.CreateReferee(referee)
	if err != nil {
		restResponse.Message = "Error to create referee: " + err.Error()
		fmt.Println(restResponse.Message)
		return restResponse
	}

	if !status {
		restResponse.Message = "Error to create referee"
		fmt.Println(restResponse.Message)
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	fmt.Println(restResponse.Message)
	return restResponse
}
