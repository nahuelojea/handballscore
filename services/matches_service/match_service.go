package matches_service

import (
	"errors"
	"time"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

type GetMatchesOptions struct {
	PhaseId       string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateMatches(associationID string, matches []models.Match) ([]string, bool, error) {
	return matches_repository.CreateMatches(associationID, matches)
}

func CreateMatch(association_id string, match models.Match) (string, bool, error) {
	return matches_repository.CreateMatch(association_id, match)
}

func GetMatch(ID string) (models.Match, bool, error) {
	return matches_repository.GetMatch(ID)
}

func GetMatches(filterOptions GetMatchesOptions) ([]models.Match, int64, error) {
	filters := matches_repository.GetMatchesOptions{
		PhaseId:       filterOptions.PhaseId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return matches_repository.GetMatchesFilteredAndPaginated(filters)
}

func ProgramMatch(matchTime time.Time, place string, id string) (bool, error) {
	if matchTime.Compare(time.Now()) < 1 {
		return false, errors.New("The date cannot be earlier than the current date")
	}
	return matches_repository.ProgramMatch(matchTime, place, id)
}

func StartMatch(startMatchRequest dto.StartMatchRequest, id string) (bool, error) {
	if len(startMatchRequest.PlayersLocal) < 1 {
		return false, errors.New("There must be a minimum of one player on the home team")
	}
	if len(startMatchRequest.PlayersVisiting) < 1 {
		return false, errors.New("There must be a minimum of one player on the visiting team")
	}
	if len(startMatchRequest.Referees) < 1 {
		return false, errors.New("There must be a minimum of one referee")
	}
	if len(startMatchRequest.Timekeeper) < 1 {
		return false, errors.New("Timekeeper is required")
	}
	if len(startMatchRequest.Scorekeeper) < 1 {
		return false, errors.New("Scorekeeper is required")
	}
	return matches_repository.StartMatch(startMatchRequest, id)
}
