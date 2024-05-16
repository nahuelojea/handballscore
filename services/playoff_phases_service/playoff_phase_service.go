package playoff_phases_service

import (
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_phases_repository"
	"github.com/nahuelojea/handballscore/services/playoff_rounds_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetPlayoffPhasesOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func GetPlayoffPhases(filterOptions GetPlayoffPhasesOptions) ([]models.PlayoffPhase, int64, int, error) {
	filters := playoff_phases_repository.GetPlayoffPhasesOptions{
		TournamentCategoryId: filterOptions.TournamentCategoryId,
		AssociationId:        filterOptions.AssociationId,
		Page:                 filterOptions.Page,
		PageSize:             filterOptions.PageSize,
		SortField:            filterOptions.SortField,
		SortOrder:            filterOptions.SortOrder,
	}
	return playoff_phases_repository.GetPlayoffPhases(filters)
}

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
