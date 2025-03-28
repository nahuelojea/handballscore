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
	err, status := generateMatchSanctions(&endMatch.Match)

	endMatch.UpdatePlayersSanctions = models.StepStatus{
		IsDone: err == nil,
		Status: func() string {
			if err != nil {
				return err.Error()
			} else {
				return status
			}
		}(),
	}

	fmt.Println("UpdatePlayerSanctionsHandler Status: ", endMatch.UpdatePlayersSanctions.Status)

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func generateMatchSanctions(match *models.Match) (error, string) {
	status := "Step ignored"

	if match.Status == models.Suspended {
		return nil, status
	}

	getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
		MatchId:       match.Id.Hex(),
		HasBlueCard:   true,
		AssociationId: match.AssociationId,
	}
	sanctionatedPlayers, _, _, err := match_players_repository.GetMatchPlayers(getPlayersOptions)

	if err != nil {
		status = "Error to get players with blue card: " + err.Error()
	} else {
		for _, player := range sanctionatedPlayers {
			playerSanction := models.PlayerSanction{
				IssueDate:      match.Date,
				Description:    match.Comments,
				SanctionStatus: models.PendingReview,
				PlayerId:       player.PlayerId,
				MatchId:        match.Id.Hex(),
			}

			_, _, err = player_sanctions_repository.CreatePlayerSanction(match.AssociationId, playerSanction)
			if err != nil {
				fmt.Println("Error to create player sanction: " + err.Error())
			}
		}
	}
	return err, status
}
