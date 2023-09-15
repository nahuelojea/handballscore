package associations_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/associations_repository"
)

type GetAssociationsOptions struct {
	Name      string
	Page      int
	PageSize  int
	SortField string
	SortOrder int
}

func GetAssociation(ID string) (models.Association, bool, error) {
	return associations_repository.GetAssociation(ID)
}

func GetAssociations(filterOptions GetAssociationsOptions) ([]models.Association, int64, error) {

	filters := associations_repository.GetAssociationsOptions{
		Name:      filterOptions.Name,
		Page:      filterOptions.Page,
		PageSize:  filterOptions.PageSize,
		SortField: filterOptions.SortField,
		SortOrder: filterOptions.SortOrder,
	}

	return associations_repository.GetAssociationsFilteredAndPaginated(filters)
}
