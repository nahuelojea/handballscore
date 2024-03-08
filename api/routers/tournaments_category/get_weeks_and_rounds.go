package tournaments

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	categories_dto "github.com/nahuelojea/handballscore/dto/tournament_categories"
	"github.com/nahuelojea/handballscore/services/league_phase_weeks_service"
	"github.com/nahuelojea/handballscore/services/league_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_rounds_service"
)

func GetWeeksAndRounds(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var weeksAndRoundsResponse []categories_dto.WeeksAndRoundsResponse

	id := request.QueryStringParameters["id"]
	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	associationId := claim.AssociationId

	if len(associationId) < 1 {
		return dto.RestResponse{
			Status:  http.StatusBadRequest,
			Message: "'associationId' is mandatory",
		}
	}
	if len(id) < 1 {
		return dto.RestResponse{
			Status:  http.StatusBadRequest,
			Message: "'id' param is mandatory",
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	filterLeaguePhase := league_phases_service.GetLeaguePhasesOptions{
		TournamentCategoryId: id,
		AssociationId:        associationId,
		Page:                 page,
		PageSize:             pageSize,
	}

	leaguePhases, _, _, err := league_phases_service.GetLeaguePhases(filterLeaguePhase)
	if err != nil {
		return dto.RestResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error to get league phase: " + err.Error(),
		}
	}

	if len(leaguePhases) > 0 {
		leaguePhase := leaguePhases[0]
		filterLeaguePhaseWeek := league_phase_weeks_service.GetLeaguePhaseWeeksOptions{
			LeaguePhaseId: leaguePhase.Id.Hex(),
			AssociationId: associationId,
			Page:          page,
			PageSize:      pageSize,
		}

		leaguePhaseWeeks, _, _, _ := league_phase_weeks_service.GetLeaguePhaseWeeks(filterLeaguePhaseWeek)

		for _, week := range leaguePhaseWeeks {
			weeksAndRounds := categories_dto.WeeksAndRoundsResponse{
				Description:       "Jornada " + strconv.Itoa(week.Number),
				LeaguePhaseWeekId: week.Id.Hex(),
			}
			weeksAndRoundsResponse = append(weeksAndRoundsResponse, weeksAndRounds)
		}
	}

	filterPlayoffPhase := playoff_phases_service.GetPlayoffPhasesOptions{
		TournamentCategoryId: id,
		AssociationId:        associationId,
		Page:                 page,
		PageSize:             pageSize,
	}

	playoffPhases, _, _, err := playoff_phases_service.GetPlayoffPhases(filterPlayoffPhase)
	if err != nil {
		return dto.RestResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error to get playoff phase: " + err.Error(),
		}
	}

	if len(playoffPhases) > 0 {
		playoffPhase := playoffPhases[0]

		filterPlayoffRound := playoff_rounds_service.GetPlayoffRoundsOptions{
			PlayoffPhaseId: playoffPhase.Id.Hex(),
			AssociationId:  associationId,
			Page:           page,
			PageSize:       pageSize,
		}

		playoffRounds, _, _, _ := playoff_rounds_service.GetPlayoffRounds(filterPlayoffRound)

		for _, round := range playoffRounds {
			weeksAndRounds := categories_dto.WeeksAndRoundsResponse{
				Description:       round.PlayoffRoundNameTraduction(),
				LeaguePhaseWeekId: round.Id.Hex(),
			}
			weeksAndRoundsResponse = append(weeksAndRoundsResponse, weeksAndRounds)
		}
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: int64(len(weeksAndRoundsResponse)),
		TotalPages:   int(math.Ceil(float64(len(weeksAndRoundsResponse)) / float64(pageSize))),
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        weeksAndRoundsResponse,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		return dto.RestResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error formatting categories to JSON: " + err.Error(),
		}
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
