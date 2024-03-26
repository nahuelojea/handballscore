package matches

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/matches_service"
	"github.com/nahuelojea/handballscore/services/playoff_round_keys_service"
)

func GetMatchesByJourney(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	leaguePhaseWeekId := request.QueryStringParameters["leaguePhaseWeekId"]
	playoffRoundId := request.QueryStringParameters["playoffRoundId"]
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

	filterPlayoffRoundKeys := playoff_round_keys_service.GetPlayoffRoundKeysOptions{
		PlayoffRoundId: playoffRoundId,
		AssociationId:  associationId,
	}
	playoffRoundKeys, _, _, _ := playoff_round_keys_service.GetPlayoffRoundKeys(filterPlayoffRoundKeys)
	var playoffRoundKeyIds []string
	if len(playoffRoundKeys) > 1 {
		for _, playoffRoundKey := range playoffRoundKeys {
			playoffRoundKeyIds = append(playoffRoundKeyIds, playoffRoundKey.Id.Hex())
		}
	}

	fmt.Println("playoffRoundKeyIds: ", playoffRoundKeyIds)

	filterMatches := matches_service.GetMatchesOptions{
		LeaguePhaseWeekId:  leaguePhaseWeekId,
		PlayoffRoundKeyIds: playoffRoundKeyIds,
		AssociationId:      associationId,
		Page:               page,
		PageSize:           pageSize,
		SortOrder:          1,
	}

	matchesList, totalRecords, totalPages, err := matches_service.GetMatchesByJourney(filterMatches)
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
