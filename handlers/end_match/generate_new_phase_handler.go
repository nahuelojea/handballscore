package end_match

import (
	"github.com/nahuelojea/handballscore/models"
)

type GenerateNewPhaseHandler struct {
	BaseEndMatchHandler
}

func (c *GenerateNewPhaseHandler) HandleEndMatch(endMatch *models.EndMatch) {

	/*var err error
	var status string = "Step ignored"

	if endMatch.CurrentPhase == models.League_Current_Phase {

		playoffPhases, _, _, err := playoff_phases_service.GetPlayoffPhases(playoff_phases_service.GetPlayoffPhasesOptions{
			TournamentCategoryId: endMatch.CurrentTournamentCategory.TournamentId,
		})
		if err != nil {
			endMatch.GenerateNewPhase = models.StepStatus{IsDone: false, Status: err.Error()}
		}

		if len(playoffPhases) != 0 {
			matches, err := matches_service.GetPendingMatchesByLeaguePhaseId(endMatch.CurrentLeaguePhase.LeaguePhase.Id.Hex())
			if err != nil {
				endMatch.GenerateNewPhase = models.StepStatus{IsDone: false, Status: err.Error()}
			}

			if len(matches) == 0 {
				playoffRoundKey, _, _, err := playoff_rounds_service.GetPlayoffRounds(playoff_rounds_service.GetPlayoffRoundsOptions{
					PlayoffPhaseId: playoffPhases[0].Id.Hex(),
			}
		}

	}

	if err != nil {
		endMatch.GenerateNewPhase = models.StepStatus{IsDone: false, Status: err.Error()}
	} else {
		endMatch.GenerateNewPhase = models.StepStatus{IsDone: true, Status: status}
	}*/

	if c.GetNext() != nil {
		c.GetNext().HandleEndMatch(endMatch)
	}
}
