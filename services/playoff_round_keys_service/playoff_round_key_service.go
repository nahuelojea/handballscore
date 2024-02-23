package playoff_round_keys_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_round_keys_repository"
)

func CreatePlayoffRound(association_id string, playoffRoundKey models.PlayoffRoundKey) (string, bool, error) {
	return playoff_round_keys_repository.CreatePlayoffRoundKey(association_id, playoffRoundKey)
}

func CreatePlayoffRoundKeys(association_id string, playoffRoundKeys []models.PlayoffRoundKey) ([]string, bool, error) {
	return playoff_round_keys_repository.CreatePlayoffRoundKeys(association_id, playoffRoundKeys)
}
