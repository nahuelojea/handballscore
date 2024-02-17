package league_phase_weeks_service

import (
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phase_weeks_repository"
)

type GetLeaguePhaseWeeksOptions struct {
	LeaguePhaseId string
	Number        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateLeaguePhaseWeeks(association_id string, leaguePhaseWeeks []models.LeaguePhaseWeek) ([]string, bool, error) {
	return league_phase_weeks_repository.CreateLeaguePhaseWeeks(association_id, leaguePhaseWeeks)
}

func GetLeaguePhaseWeek(id string) (models.LeaguePhaseWeek, bool, error) {
	return league_phase_weeks_repository.GetLeaguePhaseWeek(id)
}

func GetLeaguePhaseWeeks(filterOptions GetLeaguePhaseWeeksOptions) ([]models.LeaguePhaseWeek, int64, error) {
	filters := league_phase_weeks_repository.GetLeaguePhaseWeeksOptions{
		LeaguePhaseId: filterOptions.LeaguePhaseId,
		Number:        filterOptions.Number,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return league_phase_weeks_repository.GetLeaguePhaseWeeks(filters)
}
