package associations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/associations_service"
)

func GetAssociations(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}

	filterOptions := associations_service.GetAssociationsOptions{
		Name:      name,
		Page:      page,
		PageSize:  pageSize,
		SortField: "name",
		SortOrder: 1,
	}

	associationsList, totalRecords, totalPages, err := associations_service.GetAssociations(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get associations: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        associationsList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting associations to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
