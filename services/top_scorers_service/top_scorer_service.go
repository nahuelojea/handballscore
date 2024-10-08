package top_scorers_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/top_scorers_repository"
)

type GetTopScorersOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Name                 string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func GetTopScorers(filterOptions GetTopScorersOptions) ([]models.TopScorer, int64, int, error) {
	filters := top_scorers_repository.GetTopScorersOptions{
		TournamentCategoryId: filterOptions.TournamentCategoryId,
		AssociationId:        filterOptions.AssociationId,
		Name:                 filterOptions.Name,
		Page:                 filterOptions.Page,
		PageSize:             filterOptions.PageSize,
	}
	return top_scorers_repository.GetTopScorers(filters)
}
