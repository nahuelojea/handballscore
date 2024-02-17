package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func GetTournamentsCategory(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	categoryId := request.QueryStringParameters["category_id"]
	tournamentId := request.QueryStringParameters["tournament_id"]
	status := request.QueryStringParameters["status"]
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

	filterOptions := tournaments_service.GetTournamentsCategoryOptions{
		Name:          name,
		CategoryId:    categoryId,
		TournamentId:  tournamentId,
		Status:        status,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
	}

	tournamentsList, totalRecords, totalPages, err := tournaments_service.GetTournamentsCategory(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get tournaments category: " + err.Error()
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
		response.Message = "Error formatting tournaments category to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
