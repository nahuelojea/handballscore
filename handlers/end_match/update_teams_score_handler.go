package end_match

import (
	"errors"
	"fmt"

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
		err = updateLeagueStandings(endMatch)
	case models.Playoff_Current_Phase:
		err = updatePlayoffStandings(endMatch)
	}

	endMatch.UpdateTeamsScore = models.StepStatus{
		IsDone: err == nil,
		Status: func() string {
			if err != nil {
				return err.Error()
			}
			return "Teams scores updated"
		}(),
	}

	fmt.Println("UpdateTeamsScore Status: ", endMatch.UpdateTeamsScore.Status)

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func updateLeagueStandings(endMatch *models.EndMatch) error {
	leaguePhase := &endMatch.CurrentLeaguePhase.LeaguePhase
	if err := updateStandings(&endMatch.Match, leaguePhase.TeamsRanking[:]); err != nil {
		return err
	}
	_, err := league_phases_repository.UpdateTeamsRanking(*leaguePhase, leaguePhase.Id.Hex())
	return err
}

func updatePlayoffStandings(endMatch *models.EndMatch) error {
	playoffRoundKey := &endMatch.CurrentPlayoffPhase.PlayoffRoundKey
	if err := updateStandings(&endMatch.Match, playoffRoundKey.TeamsRanking[:]); err != nil {
		return err
	}

	matchResult := models.MatchResult{
		TeamHomeGoals: endMatch.Match.GoalsHome.Total,
		TeamAwayGoals: endMatch.Match.GoalsAway.Total,
	}

	playoffRoundKey.MatchResults = append(playoffRoundKey.MatchResults, matchResult)

	_, err := playoff_round_keys_repository.UpdatePlayoffRoundKey(*playoffRoundKey, playoffRoundKey.Id.Hex())
	return err
}

func updateStandings(match *models.Match, teamsScores []models.TeamScore) error {
	if match.GoalsHome.Total < 0 || match.GoalsAway.Total < 0 {
		return errors.New("Invalid match result, goals can't be negative")
	}

	homeTeam, awayTeam := findTeamInStandings(match.TeamHome, teamsScores), findTeamInStandings(match.TeamAway, teamsScores)
	if homeTeam == nil || awayTeam == nil {
		return errors.New("Teams not found in standings")
	}

	if match.Status != models.Suspended {
		updateTeamScores(match, homeTeam, awayTeam)
	}

	return nil
}

func updateTeamScores(match *models.Match, homeTeam, awayTeam *models.TeamScore) {
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

func findTeamInStandings(teamId models.TournamentTeamId, standings []models.TeamScore) *models.TeamScore {
	for i := range standings {
		if standings[i].TeamId == teamId {
			return &standings[i]
		}
	}
	return nil
}
