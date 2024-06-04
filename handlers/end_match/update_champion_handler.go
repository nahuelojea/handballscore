package end_match

import (
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_phases_repository"
	tournaments_repository "github.com/nahuelojea/handballscore/repositories/tournaments_category_repository"
)

type UpdateChampionHandler struct {
	BaseEndMatchHandler
}

func (c *UpdateChampionHandler) HandleEndMatch(endMatch *models.EndMatch) {
	var status string

	switch endMatch.CurrentPhase {
	case models.League_Current_Phase:
		handleLeaguePhaseChampion(endMatch, &status)
	case models.Playoff_Current_Phase:
		handlePlayoffPhaseChampion(endMatch, &status)
	}

	endMatch.GenerateNewPhase = models.StepStatus{
		IsDone: status != "",
		Status: status,
	}

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func handleLeaguePhaseChampion(endMatch *models.EndMatch, status *string) {
	leaguePhase := &endMatch.CurrentLeaguePhase.LeaguePhase

	if leaguePhase.Finished {
		filterPlayoffPhase := playoff_phases_repository.GetPlayoffPhasesOptions{
			TournamentCategoryId: endMatch.CurrentTournamentCategory.Id.Hex(),
			AssociationId:        endMatch.CurrentTournamentCategory.AssociationId,
		}

		playoffPhase, _, _, _ := playoff_phases_repository.GetPlayoffPhases(filterPlayoffPhase)
		if len(playoffPhase) != 0 {
			return
		}

		endMatch.CurrentTournamentCategory.Champion = leaguePhase.Winner
		endMatch.CurrentTournamentCategory.EndDate = time.Now()
		endMatch.CurrentTournamentCategory.Status = models.Ended

		tournaments_repository.UpdateTournamentCategory(endMatch.CurrentTournamentCategory, endMatch.CurrentTournamentCategory.Id.Hex())
		*status = "Champion updated"
	}
}

func handlePlayoffPhaseChampion(endMatch *models.EndMatch, status *string) {
	playoffRoundKey := &endMatch.CurrentPlayoffPhase.PlayoffRoundKey

	if playoffRoundKey.NextRoundKeyId != "" && playoffRoundKey.Finished {
		endMatch.CurrentTournamentCategory.Champion = playoffRoundKey.Winner
		endMatch.CurrentTournamentCategory.EndDate = time.Now()
		endMatch.CurrentTournamentCategory.Status = models.Ended

		tournaments_repository.UpdateTournamentCategory(endMatch.CurrentTournamentCategory, endMatch.CurrentTournamentCategory.Id.Hex())
		*status = "Champion updated"
	}
}
