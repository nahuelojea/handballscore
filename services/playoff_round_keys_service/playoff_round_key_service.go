package playoff_round_keys_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_round_keys_repository"
)

type GetPlayoffRoundKeysOptions struct {
	PlayoffRoundId string
	AssociationId  string
	Page           int
	PageSize       int
	SortField      string
	SortOrder      int
}

func GetPlayoffRoundKeys(filterOptions GetPlayoffRoundKeysOptions) ([]models.PlayoffRoundKey, int64, int, error) {
	filters := playoff_round_keys_repository.GetPlayoffRoundKeysOptions{
		PlayoffRoundId: filterOptions.PlayoffRoundId,
		AssociationId:  filterOptions.AssociationId,
		Page:           filterOptions.Page,
		PageSize:       filterOptions.PageSize,
		SortField:      filterOptions.SortField,
		SortOrder:      filterOptions.SortOrder,
	}
	return playoff_round_keys_repository.GetPlayoffRoundKeys(filters)
}

func CreatePlayoffRound(association_id string, playoffRoundKey models.PlayoffRoundKey) (string, bool, error) {
	return playoff_round_keys_repository.CreatePlayoffRoundKey(association_id, playoffRoundKey)
}

func CreatePlayoffRoundKeys(association_id string, playoffRoundKeys []models.PlayoffRoundKey) ([]string, bool, error) {
	return playoff_round_keys_repository.CreatePlayoffRoundKeys(association_id, playoffRoundKeys)
}
