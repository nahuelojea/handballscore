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
	"github.com/nahuelojea/handballscore/services/league_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_phases_service"
	"github.com/nahuelojea/handballscore/services/tournaments_service"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetTournamentsCategoryOptions struct {
	Name          string
	CategoryId    string
	TournamentId  string
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
		return "", false, errors.New(fmt.Sprintf("Error to get tournament: %s", err.Error()))
	}

	var tournamentCategory models.TournamentCategory

	tournamentCategory.CategoryId = category.Id.Hex()
	tournamentCategory.TournamentId = tournament.Id.Hex()
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

	switch tournamentRequest.Format {
	case TournamentDTO.League_Format:
		if !reflect.DeepEqual(tournamentRequest.LeaguePhase, TournamentDTO.LeaguePhaseRequest{}) {
			return league_phases_service.CreateTournamentLeaguePhase(tournamentCategory,
				tournamentCategoryId,
				tournamentRequest.LeaguePhase.HomeAndAway)
		} else {
			return "", false, errors.New("League data is required")
		}
	case TournamentDTO.Playoff_Format:
		if !reflect.DeepEqual(tournamentRequest.PlayoffPhase, TournamentDTO.PlayoffPhaseRequest{}) {
			return playoff_phases_service.CreateTournamentPlayoffPhase(tournamentCategory,
				tournamentCategoryId,
				tournamentRequest.PlayoffPhase.HomeAndAway,
				tournamentRequest.PlayoffPhase.RandomOrder)
		} else {
			return "", false, errors.New("Playoff data is required")
		}
	}

	return tournamentCategory.TournamentId, true, nil
}

func GetTournamentCategory(ID string) (models.TournamentCategory, bool, error) {
	return TournamentsRepository.GetTournamentCategory(ID)
}

func GetTournamentsCategory(filterOptions GetTournamentsCategoryOptions) ([]models.TournamentCategory, int64, error) {
	filters := TournamentsRepository.GetTournamentsCategoryOptions{
		Name:          filterOptions.Name,
		CategoryId:    filterOptions.CategoryId,
		TournamentId:  filterOptions.TournamentId,
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
