package coaches

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
)

func GetCoachs(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	surname := request.QueryStringParameters["surname"]
	dni := request.QueryStringParameters["dni"]
	teamId := request.QueryStringParameters["teamId"]
	associationId := claim.AssociationId

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

	filterOptions := coaches_repository.GetCoachsOptions{
		Name:          name,
		Surname:       surname,
		Dni:           dni,
		TeamId:        teamId,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortField:     "personal_data.surname",
		SortOrder:     1,
	}

	coachesList, totalRecords, err := coaches_repository.GetCoachsFilteredAndPaginated(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get coaches: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   int(totalRecords / int64(pageSize)),
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        coachesList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting coaches to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
