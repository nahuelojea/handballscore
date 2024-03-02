package playoff_phases_service

import (
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_phases_repository"
	"github.com/nahuelojea/handballscore/services/playoff_rounds_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePlayoffPhase(association_id string, playoffPhase models.PlayoffPhase) (string, bool, error) {
	return playoff_phases_repository.CreatePlayoffPhase(association_id, playoffPhase)
}

func CreateTournamentPlayoffPhase(tournamentCategory models.TournamentCategory, playoffPhaseConfig models.PlayoffPhaseConfig) (string, bool, error) {

	var playoffPhase models.PlayoffPhase

	playoffPhase.TournamentCategoryId = tournamentCategory.Id.Hex()
	playoffPhase.Config = playoffPhaseConfig

	if playoffPhaseConfig.ClassifiedNumber == 0 {
		playoffPhase.Teams = tournamentCategory.Teams
	}

	playoffPhaseIdStr, _, err := CreatePlayoffPhase(tournamentCategory.AssociationId, playoffPhase)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff phase: %s", err.Error()))
	}

	playoffPhaseId, err := primitive.ObjectIDFromHex(playoffPhaseIdStr)
	if err != nil {
		return "", false, err
	}

	playoffPhase.Id = playoffPhaseId

	playoff_rounds_service.CreateTournamentPlayoffRounds(tournamentCategory.AssociationId, playoffPhase)

	return tournamentCategory.Id.Hex(), true, nil
}
