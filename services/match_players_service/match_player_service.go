package match_players_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/config/firebase"
	dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

func CreateMatchPlayer(association_id string, matchPlayerRequest dto.MatchPlayerRequest) (string, bool, error) {
	match, _, err := matches_repository.GetMatch(matchPlayerRequest.MatchId)
	if err != nil {
		return "", false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.Programmed &&
		match.Status != models.FirstHalf &&
		match.Status != models.SecondHalf {
		return "", false, errors.New("The player cannot be added in this match instance")
	}

	matchPlayerView, _, _, err := match_players_repository.GetMatchPlayers(match_players_repository.GetMatchPlayerOptions{
		MatchId:       matchPlayerRequest.MatchId,
		PlayerId:      matchPlayerRequest.PlayerId,
		AssociationId: association_id,
	})
	if err != nil {
		return "", false, errors.New("Error to get match player: " + err.Error())
	}

	if len(matchPlayerView) > 0 {
		return "", false, errors.New("The player is already in the match")
	}

	matchPlayer := models.MatchPlayer{
		MatchId:  matchPlayerRequest.MatchId,
		PlayerId: matchPlayerRequest.PlayerId,
		Number:   matchPlayerRequest.Number,
		TeamId: models.TournamentTeamId{
			TeamId:  matchPlayerRequest.Team.Id,
			Variant: matchPlayerRequest.Team.Variant,
		},
		Goals: models.Goals{
			FirstHalf:  0,
			SecondHalf: 0},
		Sanctions: models.Sanctions{
			Exclusions: []models.Exclusion{},
			YellowCard: false,
			RedCard:    false,
			BlueCard:   false,
			Report:     ""},
	}
	return match_players_repository.CreateMatchPlayer(association_id, matchPlayer)
}

func UpdateMatchPlayer(matchPlayer models.MatchPlayer, id string) (bool, error) {
	return match_players_repository.UpdateMatchPlayer(matchPlayer, id)
}

func DeleteMatchPlayer(id string) (bool, error) {
	return match_players_repository.DeleteMatchPlayer(id)
}

func GetMatchPlayer(id string) (models.MatchPlayer, bool, error) {
	return match_players_repository.GetMatchPlayer(id)
}

type GetMatchPlayerOptions struct {
	MatchId       string
	Team          models.TournamentTeamId
	PlayerId      string
	Number        int
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchPlayers(filterOptions GetMatchPlayerOptions) ([]models.MatchPlayerView, int64, int, error) {
	filters := match_players_repository.GetMatchPlayerOptions{
		MatchId:       filterOptions.MatchId,
		Team:          filterOptions.Team,
		PlayerId:      filterOptions.PlayerId,
		Number:        filterOptions.Number,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}

	return match_players_repository.GetMatchPlayers(filters)
}

func UpdateGoal(id string, addGoal bool) (bool, error) {
	matchPlayer, match, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	if matchPlayer.RedCard {
		return false, errors.New("The player has red card")
	}
	if matchPlayer.BlueCard {
		return false, errors.New("The player has blue card")
	}

	if match.Status == models.FirstHalf {
		if addGoal {
			matchPlayer.Goals.FirstHalf++
		} else if matchPlayer.Goals.FirstHalf > 0 {
			matchPlayer.Goals.FirstHalf--
		}
	} else {
		if addGoal {
			matchPlayer.Goals.SecondHalf++
		} else if matchPlayer.Goals.SecondHalf > 0 {
			matchPlayer.Goals.SecondHalf--
		} else if matchPlayer.Goals.FirstHalf > 0 {
			matchPlayer.Goals.FirstHalf--
		}
	}

	matchPlayer.Goals.Total = matchPlayer.Goals.FirstHalf + matchPlayer.Goals.SecondHalf

	if _, err := match_players_repository.UpdateGoals(matchPlayer); err != nil {
		return false, err
	}

	if _, err := RecalculateTeamGoals(match, matchPlayer.TeamId); err != nil {
		return false, fmt.Errorf("error recalculating team goals: %w", err)
	}

	firebaseDBClient, errDb := firebase.GetFirebaseDBClient()
	if errDb != nil {
		// fmt.Printf("Error getting Firebase DB client: %v. Skipping Firebase update.\n", errDb) // Removed
	} else if firebaseDBClient == nil {
		// fmt.Println("Firebase DB client is nil but no error was reported. Skipping Firebase update.") // Removed
	} else {
		matchID := matchPlayer.MatchId
		playerID := matchPlayer.PlayerId
		teamID := matchPlayer.TeamId.TeamId.Hex()

		playerGoalData := map[string]interface{}{
			"goals_first_half":  matchPlayer.Goals.FirstHalf,
			"goals_second_half": matchPlayer.Goals.SecondHalf,
			"goals_total":       matchPlayer.Goals.Total,
		}

		teamScoresData := map[string]interface{}{
			"home_score_first_half":  match.GoalsHome.FirstHalf,
			"home_score_second_half": match.GoalsHome.SecondHalf,
			"home_score_total":       match.GoalsHome.Total,
			"away_score_first_half":  match.GoalsAway.FirstHalf,
			"away_score_second_half": match.GoalsAway.SecondHalf,
			"away_score_total":       match.GoalsAway.Total,
		}

		playerPath := fmt.Sprintf("matches/%s/teams/%s/players/%s/goals", matchID, teamID, playerID)
		if err := firebaseDBClient.NewRef(playerPath).Set(context.Background(), playerGoalData); err != nil {
			// fmt.Printf("Error sending player goal update to Firebase for match %s, player %s: %v\n", matchID, playerID, err) // Removed
		} else {
			// fmt.Printf("Player goal update sent to Firebase for match %s, player %s\n", matchID, playerID) // Removed
		}

		matchScoresPath := fmt.Sprintf("matches/%s/scores", matchID)
		if err := firebaseDBClient.NewRef(matchScoresPath).Update(context.Background(), teamScoresData); err != nil {
			// fmt.Printf("Error sending team scores update to Firebase for match %s: %v\n", matchID, err) // Removed
		} else {
			// fmt.Printf("Team scores update sent to Firebase for match %s\n", matchID) // Removed
		}
	}

	return true, nil
}

func RecalculateTeamGoals(match models.Match, team models.TournamentTeamId) (bool, error) {
	teamFirstHalfGoals := 0
	teamSecondHalfGoals := 0

	getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
		MatchId:       match.Id.Hex(),
		AssociationId: match.AssociationId,
		Team:          team,
	}

	players, _, _, err := match_players_repository.GetMatchPlayers(getPlayersOptions)
	if err != nil {
		return false, errors.New("Error to get match players: " + err.Error())
	}

	for _, player := range players {
		if player.TeamId == team {
			teamFirstHalfGoals += player.Goals.FirstHalf
			teamSecondHalfGoals += player.Goals.SecondHalf
		}
	}

	if match.TeamHome == team {
		match.GoalsHome.FirstHalf = teamFirstHalfGoals
		match.GoalsHome.SecondHalf = teamSecondHalfGoals
		match.GoalsHome.Total = teamFirstHalfGoals + teamSecondHalfGoals
	} else if match.TeamAway == team {
		match.GoalsAway.FirstHalf = teamFirstHalfGoals
		match.GoalsAway.SecondHalf = teamSecondHalfGoals
		match.GoalsAway.Total = teamFirstHalfGoals + teamSecondHalfGoals
	} else {
		return false, errors.New("The team id does not match any of the two teams in the match")
	}

	return matches_repository.UpdateGoals(match, match.Id.Hex())
}

func UpdateExclusions(id string, addExclusion bool, time string) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	if addExclusion {
		if len(matchPlayer.Sanctions.Exclusions) == 3 { 
			return false, errors.New("The player already has three exclusions")
		}
		matchPlayer.Exclusions = append(matchPlayer.Exclusions, models.Exclusion{Time: time})
	} else {
		if len(matchPlayer.Exclusions) > 0 {
			matchPlayer.Exclusions = matchPlayer.Exclusions[:len(matchPlayer.Exclusions)-1]
		}
	}

	if _, err := match_players_repository.UpdateExclusions(matchPlayer); err != nil {
		return false, err
	}

	sendSanctionUpdateToFirebase(matchPlayer)

	return true, nil
}

func UpdateYellowCard(id string, addYellowCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchPlayer.YellowCard = addYellowCard

	if _, err := match_players_repository.UpdateYellowCard(matchPlayer); err != nil {
		return false, err
	}

	sendSanctionUpdateToFirebase(matchPlayer)

	return true, nil
}

func UpdateRedCard(id string, addRedCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchPlayer.RedCard = addRedCard

	if _, err := match_players_repository.UpdateRedCard(matchPlayer); err != nil {
		return false, err
	}

	sendSanctionUpdateToFirebase(matchPlayer)

	return true, nil
}

func UpdateNumber(id string, number int) (bool, error) {
	matchPlayer, _, err := match_players_repository.GetMatchPlayer(id)
	if err != nil {
		return false, errors.New("Error to get match player: " + err.Error())
	}

	matchPlayer.Number = number

	return match_players_repository.UpdateMatchPlayer(matchPlayer, id)
}

func UpdateBlueCard(id, report string, addBlueCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchPlayer.BlueCard = addBlueCard
	if addBlueCard {
		matchPlayer.Report = report
	} else {
		matchPlayer.Report = ""
	}

	if _, err := match_players_repository.UpdateBlueCard(matchPlayer); err != nil {
		return false, err
	}

	sendSanctionUpdateToFirebase(matchPlayer)

	return true, nil
}

func getMatchPlayerAvailableToAction(id string) (models.MatchPlayer, models.Match, error) {
	matchPlayer, _, err := match_players_repository.GetMatchPlayer(id)
	if err != nil {
		return models.MatchPlayer{}, models.Match{}, errors.New("Error to get match player: " + err.Error())
	}

	match, _, err := matches_repository.GetMatch(matchPlayer.MatchId)
	if err != nil {
		return models.MatchPlayer{}, models.Match{}, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.FirstHalf && match.Status != models.SecondHalf {
		return models.MatchPlayer{}, models.Match{}, errors.New("The match must be in progress")
	}

	return matchPlayer, match, nil
}

// sendSanctionUpdateToFirebase sends the player's complete sanction data to Firebase.
func sendSanctionUpdateToFirebase(matchPlayer models.MatchPlayer) {
	firebaseDBClient, errDb := firebase.GetFirebaseDBClient()
	if errDb != nil {
		// fmt.Printf("Error getting Firebase DB client: %v. Skipping sanction update for player %s, match %s.\n", errDb, matchPlayer.PlayerId, matchPlayer.MatchId) // Removed
		return
	}
	if firebaseDBClient == nil {
		// fmt.Printf("Firebase DB client is nil but no error was reported. Skipping sanction update for player %s, match %s.\n", matchPlayer.PlayerId, matchPlayer.MatchId) // Removed
		return
	}

	matchID := matchPlayer.MatchId
	playerID := matchPlayer.PlayerId
	teamID := matchPlayer.TeamId.TeamId.Hex()

	sanctionsPath := fmt.Sprintf("matches/%s/teams/%s/players/%s/sanctions", matchID, teamID, playerID)

	if err := firebaseDBClient.NewRef(sanctionsPath).Set(context.Background(), matchPlayer.Sanctions); err != nil {
		// fmt.Printf("Error sending player (%s) sanction update to Firebase for match %s: %v\n", playerID, matchID, err) // Removed
	} else {
		// fmt.Printf("Player (%s) sanction update sent to Firebase for match %s. Sanctions: %+v\n", playerID, matchID, matchPlayer.Sanctions) // Removed
	}
}
