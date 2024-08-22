package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/top_scorers_service"
)

func GetTopScorers(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	id := request.QueryStringParameters["id"]
	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	surname := request.QueryStringParameters["surname"]
	associationId := claim.AssociationId

	if len(id) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'id' param is mandatory"
		return response
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

	filterOptions := top_scorers_service.GetTopScorersOptions{
		TournamentCategoryId: id,
		AssociationId:        associationId,
		Surname:              surname,
		Page:                 page,
		PageSize:             pageSize,
	}

	topScorersList, totalRecords, totalPages, err := top_scorers_service.GetTopScorers(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get top scorers: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        topScorersList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting top scorers to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
