package categories

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/categories_repository"
)

func GetCategory(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
	}

	category, _, err := categories_repository.GetCategory(id)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get category: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(category)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating category to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
