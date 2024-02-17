package referees

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/referees_service"
)

func GetReferees(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	surname := request.QueryStringParameters["surname"]
	dni := request.QueryStringParameters["dni"]
	gender := request.QueryStringParameters["gender"]
	onlyEnabledStr := request.QueryStringParameters["onlyEnabled"]
	associationId := claim.AssociationId

	var onlyEnabled bool
	if onlyEnabledStr != "" {
		onlyEnabled, err = strconv.ParseBool(onlyEnabledStr)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = "'onlyEnabled' param is invalid"
			return response
		}
	}

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

	filterOptions := referees_service.GetRefereesOptions{
		Name:          name,
		Surname:       surname,
		Dni:           dni,
		Gender:        gender,
		OnlyEnabled:   onlyEnabled,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortOrder:     1,
	}

	refereesList, totalRecords, totalPages, err := referees_service.GetReferees(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get referees: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
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
