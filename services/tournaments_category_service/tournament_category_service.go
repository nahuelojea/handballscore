package tournaments_service

import (
	"github.com/nahuelojea/handballscore/models"
	tournaments_repository "github.com/nahuelojea/handballscore/repositories/tournaments_category_repository"
)

type GetTournamentsCategoryOptions struct {
	Name          string
	CategoryId    string
	Status        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateTournamentCategory(association_id string, tournament models.TournamentCategory) (string, bool, error) {
	return tournaments_repository.CreateTournamentCategory(association_id, tournament)
}

func GetTournamentCategory(ID string) (models.TournamentCategory, bool, error) {
	return tournaments_repository.GetTournamentCategory(ID)
}

func GetTournamentsCategory(filterOptions GetTournamentsCategoryOptions) ([]models.TournamentCategory, int64, error) {
	filters := tournaments_repository.GetTournamentsCategoryOptions{
		Name:          filterOptions.Name,
		CategoryId:    filterOptions.CategoryId,
		Status:        filterOptions.Status,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.Name,
		SortOrder:     filterOptions.SortOrder,
	}
	return tournaments_repository.GetTournamentsCategories(filters)
}

func UpdateTournamentCategory(tournament models.TournamentCategory, ID string) (bool, error) {
	return tournaments_repository.UpdateTournamentCategory(tournament, ID)
}

func DeleteTournamentCategory(ID string) (bool, error) {
	return tournaments_repository.DeleteTournamentCategory(ID)
}
