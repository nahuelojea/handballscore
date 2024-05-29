package end_match

import (
	"errors"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phases_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_round_keys_repository"
)

type UpdateTeamsScoreHandler struct {
	BaseEndMatchHandler
}

func (c *UpdateTeamsScoreHandler) HandleEndMatch(endMatch *models.EndMatch) {
	var err error

	switch endMatch.CurrentPhase {
	case models.League_Current_Phase:
		leaguePhase := &endMatch.CurrentLeaguePhase.LeaguePhase
		err = UpdateStandings(&endMatch.Match, leaguePhase.TeamsRanking[:])
		if err == nil {
			_, err = league_phases_repository.UpdateTeamsRanking(*leaguePhase, leaguePhase.Id.Hex())
		}
	case models.Playoff_Current_Phase:
		playoffRoundKey := &endMatch.CurrentPlayoffPhase.PlayoffRoundKey
		err = UpdateStandings(&endMatch.Match, playoffRoundKey.TeamsRanking[:])
		if err == nil {
			_, err = playoff_round_keys_repository.UpdateTeamsRanking(*playoffRoundKey, playoffRoundKey.Id.Hex())
			if err == nil {
				_, err = playoff_round_keys_repository.UpdateTeamsRanking(*playoffRoundKey, playoffRoundKey.Id.Hex())
			}
		}
	}

	if err != nil {
		endMatch.UpdateTeamsScore = models.StepStatus{IsDone: false, Status: err.Error()}
	} else {
		endMatch.UpdateTeamsScore = models.StepStatus{IsDone: true, Status: "Teams scores updated"}
	}

	if c.GetNext() != nil {
		c.GetNext().HandleEndMatch(endMatch)
	}
}

func UpdateStandings(match *models.Match, teamsScores []models.TeamScore) error {
	if match.GoalsHome.Total < 0 || match.GoalsAway.Total < 0 {
		return errors.New("Invalid match result, goals can't be negative")
	}

	homeTeam := findTeamInStandings(match.TeamHome, teamsScores)
	awayTeam := findTeamInStandings(match.TeamAway, teamsScores)

	if homeTeam == nil || awayTeam == nil {
		return errors.New("Teams not found in standings")
	}

	if match.Status != models.Suspended {
		switch {
		case match.GoalsHome.Total > match.GoalsAway.Total:
			if !(match.GoalsHome.Total == 9 && match.GoalsAway.Total == 0) {
				awayTeam.Points++
			}
			homeTeam.Points += 3
			homeTeam.Wins++
			awayTeam.Losses++
		case match.GoalsHome.Total < match.GoalsAway.Total:
			if !(match.GoalsAway.Total == 9 && match.GoalsHome.Total == 0) {
				homeTeam.Points++
			}
			awayTeam.Points += 3
			awayTeam.Wins++
			homeTeam.Losses++
		default:
			homeTeam.Points += 2
			awayTeam.Points += 2
			homeTeam.Draws++
			awayTeam.Draws++
		}

		homeTeam.GoalsScored += match.GoalsHome.Total
		homeTeam.GoalsConceded += match.GoalsAway.Total
		awayTeam.GoalsScored += match.GoalsAway.Total
		awayTeam.GoalsConceded += match.GoalsHome.Total
		homeTeam.Matches++
		awayTeam.Matches++
	}

	return nil
}

func findTeamInStandings(teamId models.TournamentTeamId, standings []models.TeamScore) *models.TeamScore {
	for i := range standings {
		if standings[i].TeamId == teamId {
			return &standings[i]
		}
	}
	return nil
}
