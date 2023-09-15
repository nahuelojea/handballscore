package categories_service

import (
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/categories_repository"
)

type GetCategoriesOptions struct {
	Name          string
	Gender        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateCategory(association_id string, category models.Category) (string, bool, error) {
	return categories_repository.CreateCategory(association_id, category)
}

func GetCategory(ID string) (models.Category, bool, error) {
	return categories_repository.GetCategory(ID)
}

func GetCategories(filterOptions GetCategoriesOptions) ([]models.Category, int64, error) {
	filters := categories_repository.GetCategoriesOptions{
		Name:          filterOptions.Name,
		Gender:        filterOptions.Gender,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}

	return categories_repository.GetCategoriesFilteredAndPaginated(filters)
}

func UpdateCategory(category models.Category, ID string) (bool, error) {
	return categories_repository.UpdateCategory(category, ID)
}

func DisableCategory(ID string) (bool, error) {
	return categories_repository.DisableCategory(ID)
}

func GetLimitYearsByCategory(ID string) (int, int, string, error) {
	category, _, err := GetCategory(ID)
	if err != nil {
		return 0, 0, "", err
	}

	ageLimitFromYear := time.Now().Year() - category.AgeLimitTo
	ageLimitToYear := time.Now().Year() - category.AgeLimitFrom

	return ageLimitFromYear, ageLimitToYear, category.Gender, nil
}
