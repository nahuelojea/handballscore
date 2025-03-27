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
	var err error

	err, status := generateMatchSanctions(endMatch)

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

func generateMatchSanctions(endMatch *models.EndMatch) (error, string) {
	status := "Step ignored"

	getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
		MatchId:       endMatch.Match.Id.Hex(),
		HasBlueCard:   true,
		AssociationId: endMatch.Match.AssociationId,
	}
	sanctionatedPlayers, _, _, err := match_players_repository.GetMatchPlayers(getPlayersOptions)

	if err != nil {
		status = "Error to get players with blue card: " + err.Error()
	} else {
		for _, player := range sanctionatedPlayers {
			playerSanction := models.PlayerSanction{
				IssueDate:      endMatch.Match.Date,
				Description:    endMatch.Match.Comments,
				SanctionStatus: models.PendingReview,
				PlayerId:       player.PlayerId,
				MatchId:        endMatch.Match.Id.Hex(),
			}

			_, _, err = player_sanctions_repository.CreatePlayerSanction(endMatch.Match.AssociationId, playerSanction)
			if err != nil {
				fmt.Println("Error to create player sanction: " + err.Error())
			}
		}
	}
	return err, status
}
