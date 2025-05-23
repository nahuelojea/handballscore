package places_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/places_repository"
)

type GetPlacesOptions struct {
	Name          string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreatePlace(association_id string, place models.Place) (string, bool, error) {
	return places_repository.CreatePlace(association_id, place)
}

func GetPlace(ID string) (models.Place, bool, error) {
	return places_repository.GetPlace(ID)
}

func GetPlaces(filterOptions GetPlacesOptions) ([]models.Place, int64, int, error) {
	repoFilterOptions := places_repository.GetPlacesOptions{
		Name:          filterOptions.Name,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return places_repository.GetPlaces(repoFilterOptions)
}

func UpdatePlace(place models.Place, ID string) (bool, error) {
	return places_repository.UpdatePlace(place, ID)
}

func DeletePlace(ID string) (bool, error) {
	return places_repository.DeletePlace(ID)
}
