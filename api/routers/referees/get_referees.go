package referees

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
)

func GetReferees(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	surname := request.QueryStringParameters["surname"]
	dni := request.QueryStringParameters["dni"]
	associationId := request.QueryStringParameters["associationId"]

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}

	filterOptions := referees_repository.GetRefereesOptions{
		Name:          name,
		Surname:       surname,
		Dni:           dni,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortField:     "personal_data.surname",
		SortOrder:     1,
	}

	refereesList, totalRecords, err := referees_repository.GetRefereesFilteredAndPaginated(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get referees: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   int(totalRecords / int64(pageSize)),
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        refereesList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting referees to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
