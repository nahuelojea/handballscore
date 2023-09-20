package tournaments_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/tournaments_repository"
)

type GetTournamentsOptions struct {
	Name          string
	CategoryId    string
	Status        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateTournament(association_id string, tournament models.Tournament) (string, bool, error) {
	return tournaments_repository.CreateTournament(association_id, tournament)
}

func GetTournament(ID string) (models.Tournament, bool, error) {
	return tournaments_repository.GetTournament(ID)
}

func GetTournaments(filterOptions GetTournamentsOptions) ([]models.Tournament, int64, error) {
	filters := tournaments_repository.GetTournamentsOptions{
		Name:          filterOptions.Name,
		CategoryId:    filterOptions.CategoryId,
		Status:        filterOptions.Status,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.Name,
		SortOrder:     filterOptions.SortOrder,
	}
	return tournaments_repository.GetTournamentsFilteredAndPaginated(filters)
}

func UpdateTournament(tournament models.Tournament, ID string) (bool, error) {
	return tournaments_repository.UpdateTournament(tournament, ID)
}

func DisableTournament(ID string) (bool, error) {
	return tournaments_repository.DisableTournament(ID)
}
