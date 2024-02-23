package playoff_rounds_service

import (
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_rounds_repository"
	"github.com/nahuelojea/handballscore/services/playoff_round_keys_service"
)

func CreatePlayoffRound(association_id string, playoffRound models.PlayoffRound) (string, bool, error) {
	return playoff_rounds_repository.CreatePlayoffRound(association_id, playoffRound)
}

func CreatePlayoffRounds(association_id string, playoffRounds []models.PlayoffRound) ([]string, bool, error) {
	return playoff_rounds_repository.CreatePlayoffRounds(association_id, playoffRounds)
}

func CreateTournamentPlayoffRounds(association_id string, playoffPhase models.PlayoffPhase) (string, bool, error) {

	playoffRounds, playoffRoundKeys := models.CreatePlayoffRounds(playoffPhase)

	_, _, err := CreatePlayoffRounds(association_id, playoffRounds)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff rounds: %s", err.Error()))
	}

	_, _, err = playoff_round_keys_service.CreatePlayoffRoundKeys(association_id, playoffRoundKeys)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff round keys: %s", err.Error()))
	}

	return "", true, nil
}
