package playoff_phases_service

import (
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/league_phase_weeks_service"
	"github.com/nahuelojea/handballscore/services/matches_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTournamentPlayoffPhase(tournamentCategory models.TournamentCategory, tournamentCategoryId string, homeAndAway bool, randomOrder bool) (string, bool, error) {
	var playoffPhase models.PlayoffPhase

	playoffPhase.TournamentCategoryId = tournamentCategoryId
	playoffPhase.Config.HomeAndAway = homeAndAway
	playoffPhase.Config.RandomOrder = randomOrder

	playoffPhase.Teams = tournamentCategory.Teams

	playoffPhase.InitializeTeamScores()

	leaguePhaseIdStr, _, err := CreateLeaguePhase(tournamentCategory.AssociationId, playoffPhase)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff phase: %s", err.Error()))
	}

	leaguePhaseId, err := primitive.ObjectIDFromHex(leaguePhaseIdStr)
	if err != nil {
		return "", false, err
	}

	playoffPhase.Id = leaguePhaseId
	leaguePhaseWeeks, rounds := playoffPhase.GenerateLeaguePhaseWeeks()

	_, _, err = league_phase_weeks_service.CreateLeaguePhaseWeeks(tournamentCategory.AssociationId, leaguePhaseWeeks)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase weeks: %s", err.Error()))
	}

	filterOptions := league_phase_weeks_service.GetLeaguePhaseWeeksOptions{
		AssociationId: tournamentCategory.AssociationId,
		LeaguePhaseId: leaguePhaseId.Hex(),
	}

	leaguePhaseWeeks, _, err = league_phase_weeks_service.GetLeaguePhaseWeeks(filterOptions)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to get league phase weeks: %s", err.Error()))
	}

	matches := playoffPhase.GenerateMatches(rounds, leaguePhaseWeeks)

	_, _, err = matches_service.CreateMatches(tournamentCategory.AssociationId, matches)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase matches: %s", err.Error()))
	}

	return "", true, nil
}
