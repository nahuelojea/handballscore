package end_match

import (
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phases_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_round_keys_repository"
)

type EndPhaseHandler struct {
	BaseEndMatchHandler
}

func (c *EndPhaseHandler) HandleEndMatch(endMatch *models.EndMatch) {
	var err error
	status := "Step ignored"

	switch endMatch.CurrentPhase {
	case models.League_Current_Phase:
		err = handleLeaguePhaseEnd(endMatch, &status)
	case models.Playoff_Current_Phase:
		err = handlePlayoffPhaseEnd(endMatch, &status)
	}

	endMatch.GenerateNewPhase = models.StepStatus{
		IsDone: err == nil,
		Status: func() string {
			if err != nil {
				return err.Error()
			}
			return status
		}(),
	}

	fmt.Println("EndPhaseHandler Status: ", endMatch.GenerateNewPhase.Status)

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func handleLeaguePhaseEnd(endMatch *models.EndMatch, status *string) error {
	leaguePhase := &endMatch.CurrentLeaguePhase.LeaguePhase
	matches, err := matches_repository.GetPendingMatchesByLeaguePhaseId(leaguePhase.Id.Hex())
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		league_phases_repository.ApplyOlympicTiebreaker(leaguePhase)

		leaguePhase.Winner = leaguePhase.TeamsRanking[0].TeamId

		if _, err := league_phases_repository.FinishPhase(leaguePhase.Id.Hex(), leaguePhase.Winner); err != nil {
			return err
		}
		endMatch.CurrentLeaguePhase.LeaguePhase.Finished = true
		*status = "League phase finished"
	}
	return nil
}

func handlePlayoffPhaseEnd(endMatch *models.EndMatch, status *string) error {
	playoffRoundKey := &endMatch.CurrentPlayoffPhase.PlayoffRoundKey
	matches, err := matches_repository.GetPendingMatchesByPlayoffRoundKeyId(playoffRoundKey.Id.Hex())
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		playoffRoundKey.SortTeamsRanking()

		playoffRoundKey.Winner = playoffRoundKey.TeamsRanking[0].TeamId

		if _, err := playoff_round_keys_repository.FinishRoundKey(playoffRoundKey.Id.Hex(), playoffRoundKey.Winner); err != nil {
			return err
		}
		endMatch.CurrentPlayoffPhase.PlayoffRoundKey.Finished = true
		*status = "Round key finished"
	}

	return nil
}
