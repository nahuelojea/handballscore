package top_scorers_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/top_scorers_repository"
)

type GetTopScorersOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func GetTopScorers(filterOptions GetTopScorersOptions) ([]models.TopScorer, int64, int, error) {
	filters := top_scorers_repository.GetTopScorersOptions{
		TournamentCategoryId: filterOptions.TournamentCategoryId,
		AssociationId:        filterOptions.AssociationId,
		Page:                 filterOptions.Page,
		PageSize:             filterOptions.PageSize,
		SortField:            filterOptions.SortField,
		SortOrder:            filterOptions.SortOrder,
	}
	return top_scorers_repository.GetTopScorers(filters)
}
