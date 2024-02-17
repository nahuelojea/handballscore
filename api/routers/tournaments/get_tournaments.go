package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/tournaments_service"
)

func GetTournaments(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	onlyEnabledStr := request.QueryStringParameters["onlyEnabled"]
	associationId := claim.AssociationId

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	var onlyEnabled bool
	if onlyEnabledStr != "" {
		onlyEnabled, err = strconv.ParseBool(onlyEnabledStr)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = "'onlyEnabled' param is invalid"
			return response
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}

	filterOptions := tournaments_service.GetTournamentsOptions{
		Name:          name,
		OnlyEnabled:   onlyEnabled,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortField:     "name",
		SortOrder:     1,
	}

	tournamentsList, totalRecords, totalPages, err := tournaments_service.GetTournaments(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get tournaments: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        tournamentsList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting tournaments to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
