package matches

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func GetMatchesByTeam(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	tournamentCategoryId := request.QueryStringParameters["tournamentCategoryId"]
	teamId := request.QueryStringParameters["teamId"]
	variant := request.QueryStringParameters["variant"]
	associationId := claim.AssociationId

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	if len(teamId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'teamId' param is mandatory"
		return response
	}

	if len(tournamentCategoryId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'tournamentCategoryId' param is mandatory"
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

	filterMatches := matches_service.GetMatchesOptions{
		TournamentCategoryId: tournamentCategoryId,
		Teams: []models.TournamentTeamId{
			{
				TeamId:  teamId,
				Variant: variant,
			},
		},
		AssociationId:      associationId,
		Page:               page,
		PageSize:           pageSize,
		SortField: 		"date",
		SortOrder:          1,
	}

	matchesList, totalRecords, totalPages, err := matches_service.GetMatchesByTeam(filterMatches)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get matches: " + err.Error()
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
