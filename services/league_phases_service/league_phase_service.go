package league_phases_service

import (
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phases_repository"
	"github.com/nahuelojea/handballscore/services/league_phase_weeks_service"
	"github.com/nahuelojea/handballscore/services/matches_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLeaguePhasesOptions struct {
	TournamentId  string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateLeaguePhase(association_id string, leaguePhase models.LeaguePhase) (string, bool, error) {
	return league_phases_repository.CreateLeaguePhase(association_id, leaguePhase)
}

func GetLeaguePhase(ID string) (models.LeaguePhase, bool, error) {
	return league_phases_repository.GetLeaguePhase(ID)
}

func GetLeaguePhases(filterOptions GetLeaguePhasesOptions) ([]models.LeaguePhase, int64, int, error) {
	filters := league_phases_repository.GetLeaguePhasesOptions{
		TournamentId:  filterOptions.TournamentId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return league_phases_repository.GetLeaguePhases(filters)
}

func DeleteLeaguePhase(ID string) (bool, error) {
	return league_phases_repository.DeleteLeaguePhase(ID)
}

func CreateTournamentLeaguePhase(tournamentCategory models.TournamentCategory, tournamentCategoryId string, homeAndAway bool) (string, bool, error) {
	var leaguePhase models.LeaguePhase

	leaguePhase.TournamentCategoryId = tournamentCategoryId
	leaguePhase.Config.HomeAndAway = homeAndAway
	leaguePhase.Config.ClassifiedNumber = 1

	leaguePhase.Teams = tournamentCategory.Teams

	leaguePhase.InitializeTeamScores()

	leaguePhaseIdStr, _, err := CreateLeaguePhase(tournamentCategory.AssociationId, leaguePhase)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase: %s", err.Error()))
	}

	leaguePhaseId, err := primitive.ObjectIDFromHex(leaguePhaseIdStr)
	if err != nil {
		return "", false, err
	}

	leaguePhase.Id = leaguePhaseId
	leaguePhaseWeeks, rounds := leaguePhase.GenerateLeaguePhaseWeeks()

	_, _, err = league_phase_weeks_service.CreateLeaguePhaseWeeks(tournamentCategory.AssociationId, leaguePhaseWeeks)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase weeks: %s", err.Error()))
	}

	filterOptions := league_phase_weeks_service.GetLeaguePhaseWeeksOptions{
		AssociationId: tournamentCategory.AssociationId,
		LeaguePhaseId: leaguePhaseId.Hex(),
	}

	leaguePhaseWeeks, _, _, err = league_phase_weeks_service.GetLeaguePhaseWeeks(filterOptions)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to get league phase weeks: %s", err.Error()))
	}

	matches := leaguePhase.GenerateMatches(rounds, leaguePhaseWeeks)

	_, _, err = matches_service.CreateMatches(tournamentCategory.AssociationId, matches)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase matches: %s", err.Error()))
	}

	return "", true, nil
}
