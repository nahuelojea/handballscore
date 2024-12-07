package end_match

import (
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

type LoadDataNextMatchesHandler struct {
	BaseEndMatchHandler
}

func (c *LoadDataNextMatchesHandler) HandleEndMatch(endMatch *models.EndMatch) {
	var err error

	teamHomeMatches, _, _ := matches_repository.GetMatchesByStatus(endMatch.Match.TeamHome, endMatch.CurrentTournamentCategory.Id.Hex(), models.Programmed)

	if len(teamHomeMatches) > 0 {
		for _, match := range teamHomeMatches {
			getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
				MatchId:       match.Id.Hex(),
				Team:          endMatch.Match.TeamHome,
				AssociationId: match.AssociationId,
			}
		
			players, _, _, _ := match_players_repository.GetMatchPlayers(getPlayersOptions)
			if len(players) != 0 {
				continue
			}

			err = laodMatchPlayersAndCoachesFromLastMatch(match, endMatch.Match, endMatch.Match.TeamHome)
			if err != nil {
				break
			}
		}
	}

	teamAwayMatches, _, _ := matches_repository.GetMatchesByStatus(endMatch.Match.TeamAway, endMatch.CurrentTournamentCategory.Id.Hex(), models.Programmed)

	if len(teamAwayMatches) > 0 {
		for _, match := range teamAwayMatches {
			getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
				MatchId:       match.Id.Hex(),
				Team:          endMatch.Match.TeamAway,
				AssociationId: match.AssociationId,
			}
		
			players, _, _, _ := match_players_repository.GetMatchPlayers(getPlayersOptions)
			if len(players) != 0 {
				continue
			}

			err = laodMatchPlayersAndCoachesFromLastMatch(match, endMatch.Match, endMatch.Match.TeamAway)
			if err != nil {
				break
			}
		}
	}

	endMatch.LoadDataNextMatches = models.StepStatus{
		IsDone: err == nil,
		Status: func() string {
			if err != nil {
				return err.Error()
			}
			return "Next matches data loaded"
		}(),
	}

	fmt.Println("LoadDataNextMatches Status: ", endMatch.LoadDataNextMatches.Status)

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func laodMatchPlayersAndCoachesFromLastMatch(match, currentMatch models.Match, teamId models.TournamentTeamId) error {

	if match.TeamHome == teamId {
		err := processTeamPlayersAndCoaches(match, currentMatch, match.TeamHome)
		if err != nil {
			return err
		}
	}	

	if match.TeamAway == teamId {
		err := processTeamPlayersAndCoaches(match, currentMatch, match.TeamAway)
		if err != nil {
			return err
		}
	}
	return nil
}

func processTeamPlayersAndCoaches(match, currentMatch models.Match, team models.TournamentTeamId) error {

	if time.Now().Year() != currentMatch.Date.Year() {
		return nil
	}

	getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
		MatchId:       currentMatch.Id.Hex(),
		Team:          team,
		AssociationId: match.AssociationId,
	}

	players, _, _, _ := match_players_repository.GetMatchPlayers(getPlayersOptions)
	if len(players) != 0 {

		for _, player := range players {
			matchPlayer := models.MatchPlayer{
				MatchId:  match.Id.Hex(),
				PlayerId: player.PlayerId,
				Number:   player.Number,
				TeamId: models.TournamentTeamId{
					TeamId:  player.TeamId.TeamId,
					Variant: player.TeamId.Variant,
				},
				Goals: models.Goals{
					FirstHalf:  0,
					SecondHalf: 0,
				},
				Sanctions: models.Sanctions{
					Exclusions: []models.Exclusion{},
					YellowCard: false,
					RedCard:    false,
					BlueCard:   false,
					Report:     "",
				},
			}
			_, _, err := match_players_repository.CreateMatchPlayer(match.AssociationId, matchPlayer)
			if err != nil {
				fmt.Println("Error to create match player: %s", err.Error())
			}
		}
	}

	getCoachesOptions := match_coaches_repository.GetMatchCoachOptions{
		MatchId:       currentMatch.Id.Hex(),
		Team:          team,
		AssociationId: match.AssociationId,
	}

	coaches, _, _, _ := match_coaches_repository.GetMatchCoaches(getCoachesOptions)
	if len(coaches) != 0 {
		for _, coach := range coaches {
			matchCoach := models.MatchCoach{
				MatchId: match.Id.Hex(),
				CoachId: coach.CoachId,
				TeamId: models.TournamentTeamId{
					TeamId:  coach.TeamId.TeamId,
					Variant: coach.TeamId.Variant,
				},
				Sanctions: models.Sanctions{
					Exclusions: []models.Exclusion{},
					YellowCard: false,
					RedCard:    false,
					BlueCard:   false,
					Report:     "",
				},
			}

			_, _, err := match_coaches_repository.CreateMatchCoach(match.AssociationId, matchCoach)
			if err != nil {
				fmt.Println("Error to create match coach: %s", err.Error())
			}
		}
	}

	return nil
}
