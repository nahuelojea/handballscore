package player_sanctions_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/player_sanctions_repository"
)

func CreatePlayerSanction(associationId string, playerSanction models.PlayerSanction) (string, bool, error) {
	return player_sanctions_repository.CreatePlayerSanction(associationId, playerSanction)
}

func GetPlayerSanction(id string) (models.PlayerSanction, bool, error) {
	return player_sanctions_repository.GetPlayerSanction(id)
}

func UpdatePlayerSanction(playerSanction models.PlayerSanction, id string) (bool, error) {
	return player_sanctions_repository.UpdatePlayerSanction(playerSanction, id)
}

func AddServedMatch(matchId, id string) (bool, error) {
	return player_sanctions_repository.AddServedMatch(matchId, id)
}
