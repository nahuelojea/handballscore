package tournaments_service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"

	TournamentCategoryDTO "github.com/nahuelojea/handballscore/dto/tournament_categories"
	TournamentDTO "github.com/nahuelojea/handballscore/dto/tournaments"
	"github.com/nahuelojea/handballscore/models"
	TournamentsRepository "github.com/nahuelojea/handballscore/repositories/tournaments_category_repository"
	"github.com/nahuelojea/handballscore/services/categories_service"
	"github.com/nahuelojea/handballscore/services/league_phase_weeks_service"
	"github.com/nahuelojea/handballscore/services/league_phases_service"
	"github.com/nahuelojea/handballscore/services/league_playoff_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_rounds_service"
	"github.com/nahuelojea/handballscore/services/tournaments_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetTournamentsCategoryOptions struct {
	Name          string
	CategoryId    string
	TournamentId  string
	Status        string
	ChampionId    string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateTournamentCategory(ctx context.Context, associationId string, tournamentRequest TournamentDTO.CreateTournamentCategoryRequest) (string, bool, error) {
	category, _, err := categories_service.GetCategory(tournamentRequest.CategoryId)
	if err != nil {
		return "", false, err
	}

	tournament, _, err := tournaments_service.GetTournament(tournamentRequest.TournamentId)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to get tournament: %s", err.Error()))
	}

	var tournamentCategory models.TournamentCategory

	tournamentCategory.CategoryId = category.Id.Hex()
	tournamentCategory.TournamentId = tournament.Id.Hex()
	tournamentCategory.AssociationId = associationId
	tournamentCategory.StartDate = tournamentRequest.StartDate
	tournamentCategory.Teams = assignVariants(tournamentRequest.Teams)

	tournamentCategory.Status = models.Started

	var categoryGender string
	if category.Gender == models.Female {
		categoryGender = "Femenino"
	} else {
		categoryGender = "Masculino"
	}
	tournamentCategory.Name = cases.Title(language.Spanish).String(fmt.Sprintf("%s %s %s", tournament.Name, category.Name, categoryGender))

	tournamentCategoryId, _, err := TournamentsRepository.CreateTournamentCategory(associationId, tournamentCategory)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create tournament category: %s", err.Error()))
	}

	tournamentCategory.Id, _ = primitive.ObjectIDFromHex(tournamentCategoryId)

	switch tournamentRequest.Format {
	case TournamentDTO.League_Format:
		if reflect.DeepEqual(tournamentRequest.LeaguePhase, TournamentDTO.LeaguePhaseRequest{}) {
			return "", false, errors.New("League data is required")
		}

		var leaguePhaseConfig models.LeaguePhaseConfig
		leaguePhaseConfig.HomeAndAway = tournamentRequest.LeaguePhase.HomeAndAway
		leaguePhaseConfig.ClassifiedNumber = 0

		return league_phases_service.CreateTournamentLeaguePhase(tournamentCategory, leaguePhaseConfig)

	case TournamentDTO.Playoff_Format:
		if reflect.DeepEqual(tournamentRequest.PlayoffPhase, TournamentDTO.PlayoffPhaseRequest{}) {
			return "", false, errors.New("Playoff data is required")
		}

		var playoffPhaseConfig models.PlayoffPhaseConfig
		playoffPhaseConfig.HomeAndAway = tournamentRequest.PlayoffPhase.HomeAndAway
		playoffPhaseConfig.RandomOrder = tournamentRequest.PlayoffPhase.RandomOrder
		playoffPhaseConfig.SingleMatchFinal = tournamentRequest.PlayoffPhase.SingleMatchFinal
		playoffPhaseConfig.ClassifiedNumber = 0

		return playoff_phases_service.CreateTournamentPlayoffPhase(tournamentCategory, playoffPhaseConfig)

	case TournamentDTO.League_And_Playoff_Format:
		if reflect.DeepEqual(tournamentRequest.LeagueAndPlayoff, TournamentDTO.LeagueAndPlayoffRequest{}) {
			return "", false, errors.New("League Playoff data is required")
		}

		var leaguePhaseConfig models.LeaguePhaseConfig
		leaguePhaseConfig.HomeAndAway = tournamentRequest.LeaguePhase.HomeAndAway
		leaguePhaseConfig.ClassifiedNumber = tournamentRequest.LeagueAndPlayoff.LeaguePhase.ClassifiedNumber

		var playoffPhaseConfig models.PlayoffPhaseConfig
		playoffPhaseConfig.HomeAndAway = tournamentRequest.PlayoffPhase.HomeAndAway
		playoffPhaseConfig.RandomOrder = tournamentRequest.PlayoffPhase.RandomOrder
		playoffPhaseConfig.SingleMatchFinal = tournamentRequest.PlayoffPhase.SingleMatchFinal
		playoffPhaseConfig.ClassifiedNumber = leaguePhaseConfig.ClassifiedNumber

		return league_playoff_phases_service.CreateTournamentLeagueAndPlayoffPhases(tournamentCategory, leaguePhaseConfig, playoffPhaseConfig)
	}

	return tournamentCategory.TournamentId, true, nil
}

func GetTournamentCategory(ID string) (models.TournamentCategory, bool, error) {
	return TournamentsRepository.GetTournamentCategory(ID)
}

func GetTournamentsCategory(filterOptions GetTournamentsCategoryOptions) ([]models.TournamentCategory, int64, int, error) {
	filters := TournamentsRepository.GetTournamentsCategoryOptions{
		Name:          filterOptions.Name,
		CategoryId:    filterOptions.CategoryId,
		TournamentId:  filterOptions.TournamentId,
		Status:        filterOptions.Status,
		ChampionId:    filterOptions.ChampionId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.Name,
		SortOrder:     filterOptions.SortOrder,
	}
	return TournamentsRepository.GetTournamentsCategories(filters)
}

func UpdateTournamentCategory(tournament models.TournamentCategory, ID string) (bool, error) {
	return TournamentsRepository.UpdateTournamentCategory(tournament, ID)
}

func DeleteTournamentCategory(ID string) (bool, error) {
	return TournamentsRepository.DeleteTournamentCategory(ID)
}

func assignVariants(teamIds []string) []models.TournamentTeamId {
	frequencyMap := make(map[string]int)
	teamCounter := make(map[string]int)
	var tournamentTeams []models.TournamentTeamId

	for _, id := range teamIds {
		frequencyMap[id]++
	}

	for _, id := range teamIds {
		teamCounter[id]++
		frequency := frequencyMap[id]
		counter := teamCounter[id]
		var variant string

		if frequency > 1 {
			variant = string('A' - 1 + rune(counter))
		}

		tournamentTeam := models.TournamentTeamId{
			TeamId:  id,
			Variant: variant,
		}
		tournamentTeams = append(tournamentTeams, tournamentTeam)
	}

	return tournamentTeams
}

func GetWeeksAndRounds(id string, associationId string, page int, pageSize int) ([]TournamentCategoryDTO.WeeksAndRoundsResponse, int64, int, error) {
	var weeksAndRoundsResponse []TournamentCategoryDTO.WeeksAndRoundsResponse

	filterLeaguePhase := league_phases_service.GetLeaguePhasesOptions{
		TournamentCategoryId: id,
		AssociationId:        associationId,
		Page:                 page,
		PageSize:             pageSize,
	}

	leaguePhases, _, _, err := league_phases_service.GetLeaguePhases(filterLeaguePhase)
	if err != nil {
		return nil, 0, 0, errors.New("Error to get league phase: " + err.Error())
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
			weeksAndRounds := TournamentCategoryDTO.WeeksAndRoundsResponse{
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
		return nil, 0, 0, errors.New("Error to get playoff phase: " + err.Error())
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
			weeksAndRounds := TournamentCategoryDTO.WeeksAndRoundsResponse{
				Description:    round.PlayoffRoundNameTraduction(),
				PlayoffRoundId: round.Id.Hex(),
			}
			weeksAndRoundsResponse = append(weeksAndRoundsResponse, weeksAndRounds)
		}
	}

	totalRecords := int64(len(weeksAndRoundsResponse))

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return weeksAndRoundsResponse, totalRecords, totalPages, nil
}
