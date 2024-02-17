package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/categories_service"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func GetCategoriesByTournament(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	tournamentId := request.QueryStringParameters["tournamentId"]
	associationId := claim.AssociationId

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' is mandatory"
		return response
	}

	if len(tournamentId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'tournamentId' param is mandatory"
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
		TournamentId:  tournamentId,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortOrder:     1,
	}

	tournamentsList, _, err := tournaments_service.GetTournamentsCategory(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get tournaments category: " + err.Error()
		return response
	}

	categoriesIds := GetUniqueCategoryIds(tournamentsList)

	categoriesList, totalRecords, err := categories_service.GetCategoriesByIds(categoriesIds)

	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get categories: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   int(totalRecords / int64(pageSize)),
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        categoriesList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting categories to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}

func GetUniqueCategoryIds(tournamentsList []models.TournamentCategory) []string {
	categoryIds := make([]string, 0)
	seen := make(map[string]bool)

	for _, tournament := range tournamentsList {
		if !seen[tournament.CategoryId] {
			categoryIds = append(categoryIds, tournament.CategoryId)
			seen[tournament.CategoryId] = true
		}
	}

	return categoryIds
}
