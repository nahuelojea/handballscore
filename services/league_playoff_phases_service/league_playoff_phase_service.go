package league_playoff_phases_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/league_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_phases_service"
)

func CreateTournamentLeagueAndPlayoffPhases(tournamentCategory models.TournamentCategory, leaguePhaseConfig models.LeaguePhaseConfig, playoffPhaseConfig models.PlayoffPhaseConfig) (string, bool, error) {

	_, _, err := league_phases_service.CreateTournamentLeaguePhase(tournamentCategory, leaguePhaseConfig)
	if err != nil {
		return "", false, err
	}

	_, _, err = playoff_phases_service.CreateTournamentPlayoffPhase(tournamentCategory, playoffPhaseConfig)
	if err != nil {
		return "", false, err
	}

	return tournamentCategory.Id.Hex(), false, nil
}
