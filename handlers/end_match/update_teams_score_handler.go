package end_match

import (
	"errors"

	"github.com/nahuelojea/handballscore/models"
)

type UpdateTeamsScoreHandler struct {
	BaseEndMatchHandler
}

func (c *UpdateTeamsScoreHandler) HandleEndMatch(endMatch *models.EndMatch) {
	var err error

	switch {
	case endMatch.CurrentPhase == models.League_Current_Phase:
		/*leaguePhase := endMatch.CurrentLeaguePhase.LeaguePhase
		err = UpdateStandings(&endMatch.Match, &leaguePhase.TeamsRanking)*/
	case endMatch.CurrentPhase == models.Playoff_Current_Phase:
		playoffRoundKey := endMatch.CurrentPlayoffPhase.PlayoffRoundKey
		err = UpdateStandings(&endMatch.Match, &playoffRoundKey.TeamsRanking)
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

func UpdateStandings(match *models.Match, teamsScores *[2]models.TeamScore) error {
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

func findTeamInStandings(teamId models.TournamentTeamId, standings *[2]models.TeamScore) *models.TeamScore {
	for _, team := range *standings {
		if team.TeamId == teamId {
			return &team
		}
	}
	return nil
}
