package tournaments_service

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	TournamentDTO "github.com/nahuelojea/handballscore/dto/tournaments"
	"github.com/nahuelojea/handballscore/models"
	TournamentsRepository "github.com/nahuelojea/handballscore/repositories/tournaments_category_repository"
	"github.com/nahuelojea/handballscore/services/categories_service"
	"github.com/nahuelojea/handballscore/services/league_phase_weeks_service"
	"github.com/nahuelojea/handballscore/services/league_phases_service"
	"github.com/nahuelojea/handballscore/services/matches_service"
	"github.com/nahuelojea/handballscore/services/tournaments_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetTournamentsCategoryOptions struct {
	Name          string
	CategoryId    string
	Status        string
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
		return "", false, err
	}

	var tournamentCategory models.TournamentCategory

	tournamentCategory.CategoryId = category.Id.Hex()
	tournamentCategory.TournamentId = tournament.Id.Hex()
	tournamentCategory.StartDate = tournamentRequest.StartDate
	tournamentCategory.Teams = tournamentRequest.Teams
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

	switch tournamentRequest.Format {
	case TournamentDTO.League_Format:
		return createTournamentWithLeagueFormat(tournamentRequest, tournamentCategoryId, associationId)
	}

	return tournamentCategory.TournamentId, true, nil
}

func createTournamentWithLeagueFormat(tournamentRequest TournamentDTO.CreateTournamentCategoryRequest, tournamentCategoryId, associationId string) (string, bool, error) {
	if !reflect.DeepEqual(tournamentRequest.LeaguePhase, TournamentDTO.LeaguePhaseRequest{}) {
		var leaguePhase models.LeaguePhase

		leaguePhase.TournamentCategoryId = tournamentCategoryId
		leaguePhase.HomeAndAway = tournamentRequest.LeaguePhase.HomeAndAway
		leaguePhase.ClassifiedNumber = 1
		leaguePhase.Teams = tournamentRequest.Teams
		leaguePhase.InitializeTeamScores()

		leaguePhaseIdStr, _, err := league_phases_service.CreateLeaguePhase(associationId, leaguePhase)
		if err != nil {
			return "", false, errors.New(fmt.Sprintf("Error to create league phase: %s", err.Error()))
		}

		leaguePhaseId, err := primitive.ObjectIDFromHex(leaguePhaseIdStr)
		if err != nil {
			return "", false, err
		}

		leaguePhase.Id = leaguePhaseId
		leaguePhaseWeeks := leaguePhase.GenerateLeaguePhaseWeeks()

		_, _, err = league_phase_weeks_service.CreateLeaguePhaseWeeks(associationId, leaguePhaseWeeks)
		if err != nil {
			return "", false, errors.New(fmt.Sprintf("Error to create league phase weeks: %s", err.Error()))
		}

		filterOptions := league_phase_weeks_service.GetLeaguePhaseWeeksOptions{
			AssociationId: associationId,
			LeaguePhaseId: leaguePhaseId.Hex(),
		}

		leaguePhaseWeeks, _, err = league_phase_weeks_service.GetLeaguePhaseWeeks(filterOptions)

		matches := leaguePhase.GenerateMatches(leaguePhaseWeeks)

		_, _, err = matches_service.CreateMatches(associationId, matches)
		if err != nil {
			return "", false, errors.New(fmt.Sprintf("Error to create league phase matches: %s", err.Error()))
		}
	} else {
		return "", false, errors.New("League data is required")
	}
	return "", true, nil
}

func GetTournamentCategory(ID string) (models.TournamentCategory, bool, error) {
	return TournamentsRepository.GetTournamentCategory(ID)
}

func GetTournamentsCategory(filterOptions GetTournamentsCategoryOptions) ([]models.TournamentCategory, int64, error) {
	filters := TournamentsRepository.GetTournamentsCategoryOptions{
		Name:          filterOptions.Name,
		CategoryId:    filterOptions.CategoryId,
		Status:        filterOptions.Status,
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
