package matches

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func GetMatches(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	leaguePhaseWeekId := request.QueryStringParameters["leaguePhaseWeekId"]
	playoffRoundKeyIdStr := request.QueryStringParameters["playoffRoundKeyIds"]
	//dateStr := request.QueryStringParameters["date"]
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

	var playoffRoundKeyIds []string
	if playoffRoundKeyIdStr != "" {
		playoffRoundKeyIds = strings.Split(playoffRoundKeyIdStr, ",")
	}

	filterOptions := matches_service.GetMatchesOptions{
		LeaguePhaseWeekId:  leaguePhaseWeekId,
		PlayoffRoundKeyIds: playoffRoundKeyIds,
		AssociationId:      associationId,
		Page:               page,
		PageSize:           pageSize,
		SortOrder:          1,
	}

	matchesList, totalRecords, totalPages, err := matches_service.GetMatches(filterOptions)
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
