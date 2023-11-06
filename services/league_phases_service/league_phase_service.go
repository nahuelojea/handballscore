package league_phases_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phases_repository"
)

type GetLeaguePhasesOptions struct {
	TournamentId  string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateLeaguePhase(association_id string, leaguePhase models.LeaguePhase) (string, bool, error) {
	return league_phases_repository.CreateLeaguePhase(association_id, leaguePhase)
}

func GetLeaguePhase(ID string) (models.LeaguePhase, bool, error) {
	return league_phases_repository.GetLeaguePhase(ID)
}

func GetLeaguePhases(filterOptions GetLeaguePhasesOptions) ([]models.LeaguePhase, int64, error) {
	filters := league_phases_repository.GetLeaguePhasesOptions{
		TournamentId:  filterOptions.TournamentId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return league_phases_repository.GetLeaguePhases(filters)
}

func DeleteLeaguePhase(ID string) (bool, error) {
	return league_phases_repository.DeleteLeaguePhase(ID)
}
