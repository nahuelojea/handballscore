package matches

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func GetMatchesToday(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	dateStr := request.QueryStringParameters["date"]
	exactDateStr := request.QueryStringParameters["exactDate"]
	associationId := claim.AssociationId

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	if len(dateStr) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'date' param is mandatory"
		return response
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = "Error to convert date: " + err.Error()
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

	exactDate := true
	if len(exactDateStr) > 0 {
		exactDate, err = strconv.ParseBool(exactDateStr)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = "Error to convert exactDate: " + err.Error()
			return response
		}
	}

	filterOptions := matches_service.GetMatchesOptions{
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		Date:          date,
		SortOrder:     1,
	}

	matchesList, totalRecords, totalPages, err := matches_service.GetMatchesToday(filterOptions, exactDate)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get match headers: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        matchesList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting matches to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
