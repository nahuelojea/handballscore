package tournaments_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/tournaments_repository"
)

type GetTournamentsOptions struct {
	Name          string
	OnlyEnabled   bool
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

func GetTournaments(filterOptions GetTournamentsOptions) ([]models.Tournament, int64, int, error) {
	filters := tournaments_repository.GetTournamentsOptions{
		Name:          filterOptions.Name,
		OnlyEnabled:   filterOptions.OnlyEnabled,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.Name,
		SortOrder:     filterOptions.SortOrder,
	}
	return tournaments_repository.GetTournaments(filters)
}

func UpdateTournament(tournament models.Tournament, ID string) (bool, error) {
	return tournaments_repository.UpdateTournament(tournament, ID)
}

func DeleteTournament(ID string) (bool, error) {
	return tournaments_repository.DeleteTournament(ID)
}
