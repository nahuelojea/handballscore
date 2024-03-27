package match_players

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	matches_dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/match_players_service"
)

func GetMatchPlayers(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var getMatchPlayersRequest matches_dto.GetMatchPlayersRequest
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
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

	body := ctx.Value(dto.Key("body")).(string)
	err = json.Unmarshal([]byte(body), &getMatchPlayersRequest)
	if err != nil {
		response.Message = err.Error()
		return response
	}

	var tournamentTeamId models.TournamentTeamId = models.TournamentTeamId{
		TeamId:  getMatchPlayersRequest.TournamentTeamId.TeamId,
		Variant: getMatchPlayersRequest.TournamentTeamId.Variant,
	}

	filterOptions := match_players_service.GetMatchPlayerOptions{
		MatchId:       getMatchPlayersRequest.MatchId,
		Team:          tournamentTeamId,
		PlayerId:      getMatchPlayersRequest.PlayerId,
		Number:        getMatchPlayersRequest.Number,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortOrder:     1,
	}

	matchPlayersList, totalRecords, totalPages, err := match_players_service.GetMatchPlayers(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get match players: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        matchPlayersList,
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
