package end_match

import (
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/player_sanctions_repository"
)

type UpdatePlayerSanctionsHandler struct {
	BaseEndMatchHandler
}

func (c *UpdatePlayerSanctionsHandler) HandleEndMatch(endMatch *models.EndMatch) {
	matchPlayers, err := getMatchPlayers(endMatch)
	if err != nil {
		endMatch.UpdatePlayersSanctions = models.StepStatus{
			IsDone: false,
			Status: err.Error(),
		}
		fmt.Println("UpdatePlayerSanctionsHandler Status: ", endMatch.UpdatePlayersSanctions.Status)
		return
	}

	err = processPlayerSanctions(endMatch, matchPlayers)
	if err != nil {
		endMatch.UpdatePlayersSanctions = models.StepStatus{
			IsDone: false,
			Status: err.Error(),
		}
		fmt.Println("UpdatePlayerSanctionsHandler Status: ", endMatch.UpdatePlayersSanctions.Status)
		return
	}

	err, status := generateMatchSanctions(&endMatch.Match)
	endMatch.UpdatePlayersSanctions = models.StepStatus{
		IsDone: err == nil,
		Status: func() string {
			if err != nil {
				return err.Error()
			}
			return status
		}(),
	}

	fmt.Println("UpdatePlayerSanctionsHandler Status: ", endMatch.UpdatePlayersSanctions.Status)

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func getMatchPlayers(endMatch *models.EndMatch) ([]models.MatchPlayerView, error) {
	matchPlayers, _, _, err := match_players_repository.GetMatchPlayers(match_players_repository.GetMatchPlayerOptions{
		MatchId:       endMatch.Match.Id.Hex(),
		AssociationId: endMatch.Match.AssociationId,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching match players: %w", err)
	}
	return matchPlayers, nil
}

func processPlayerSanctions(endMatch *models.EndMatch, matchPlayers []models.MatchPlayerView) error {
	matchPlayerIds := extractMatchPlayerIds(matchPlayers)

	playersSanctions, _, _, err := player_sanctions_repository.GetPlayerSanctions(player_sanctions_repository.GetPlayerSanctionsOptions{
		PlayerIds:      matchPlayerIds,
		AssociationId:  endMatch.Match.AssociationId,
		IncompleteOnly: true,
	})
	if err != nil {
		return fmt.Errorf("error fetching player sanctions: %w", err)
	}

	for _, playerSanction := range playersSanctions {
		err := updatePlayerSanction(playerSanction, endMatch.Match.Id.Hex())
		if err != nil {
			return fmt.Errorf("error updating player sanction: %w", err)
		}
	}
	return nil
}

func extractMatchPlayerIds(matchPlayers []models.MatchPlayerView) []string {
	var matchPlayerIds []string
	for _, matchPlayer := range matchPlayers {
		matchPlayerIds = append(matchPlayerIds, matchPlayer.Id.Hex())
	}
	return matchPlayerIds
}

func updatePlayerSanction(playerSanction models.PlayerSanction, matchId string) error {
	playerSanction.ServedMatches = append(playerSanction.ServedMatches, matchId)
	if len(playerSanction.ServedMatches) >= playerSanction.MatchSuspensions {
		playerSanction.SanctionStatus = models.Completed
	}

	_, err := player_sanctions_repository.UpdatePlayerSanction(playerSanction, playerSanction.Id.Hex())
	if err != nil {
		return err
	}
	return nil
}

func generateMatchSanctions(match *models.Match) (error, string) {
	if match.Status == models.Suspended {
		return nil, "Step ignored"
	}

	getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
		MatchId:       match.Id.Hex(),
		HasBlueCard:   true,
		AssociationId: match.AssociationId,
	}
	sanctionatedPlayers, _, _, err := match_players_repository.GetMatchPlayers(getPlayersOptions)
	if err != nil {
		return err, "Error fetching players with blue card: " + err.Error()
	}

	for _, player := range sanctionatedPlayers {
		err := createPlayerSanction(match, player)
		if err != nil {
			fmt.Println("Error creating player sanction: " + err.Error())
		}
	}
	return nil, "Sanctions generated successfully"
}

func createPlayerSanction(match *models.Match, player models.MatchPlayerView) error {
	playerSanction := models.PlayerSanction{
		IssueDate:      match.Date,
		SanctionStatus: models.PendingReview,
		PlayerId:       player.PlayerId,
		MatchId:        match.Id.Hex(),
	}

	_, _, err := player_sanctions_repository.CreatePlayerSanction(match.AssociationId, playerSanction)
	return err
}
