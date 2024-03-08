package matches_service

import (
	"errors"
	"time"

	dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

type GetMatchesOptions struct {
	LeaguePhaseWeekId string
	PlayoffRoundKeyId string
	AssociationId     string
	Page              int
	PageSize          int
	SortField         string
	SortOrder         int
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

func GetMatches(filterOptions GetMatchesOptions) ([]models.Match, int64, int, error) {
	filters := matches_repository.GetMatchesOptions{
		LeaguePhaseWeekId: filterOptions.LeaguePhaseWeekId,
		AssociationId:     filterOptions.AssociationId,
		Page:              filterOptions.Page,
		PageSize:          filterOptions.PageSize,
		SortField:         filterOptions.SortField,
		SortOrder:         filterOptions.SortOrder,
	}
	return matches_repository.GetMatches(filters)
}

func ProgramMatch(matchTime time.Time, place string, id string) (bool, error) {
	if matchTime.Compare(time.Now()) < 1 {
		return false, errors.New("The date cannot be earlier than the current date")
	}
	return matches_repository.ProgramMatch(matchTime, place, id)
}

func StartMatch(startMatchRequest dto.StartMatchRequest, id string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	for _, playerHome := range startMatchRequest.PlayersHome {
		matchPlayer := models.MatchPlayer{
			PlayerId: playerHome.PlayerId,
			Number:   playerHome.Number,
			TeamId:   match.TeamHome,
			Goals: models.Goals{
				FirstHalf:  0,
				SecondHalf: 0},
			Sanctions: models.Sanctions{
				Exclusions: []models.Exclusions{},
				YellowCard: false,
				RedCard:    false,
				BlueCard:   false,
				Report:     ""},
		}
		match_players_repository.CreateMatchPlayer(match.AssociationId, matchPlayer)
	}

	for _, playerAway := range startMatchRequest.PlayersAway {
		matchPlayer := models.MatchPlayer{
			PlayerId: playerAway.PlayerId,
			Number:   playerAway.Number,
			Goals: models.Goals{
				FirstHalf:  0,
				SecondHalf: 0},
			Sanctions: models.Sanctions{
				Exclusions: []models.Exclusions{},
				YellowCard: false,
				RedCard:    false,
				BlueCard:   false,
				Report:     ""},
		}
		match_players_repository.CreateMatchPlayer(match.AssociationId, matchPlayer)
	}

	for _, coachHome := range startMatchRequest.CoachsHome {
		matchCoach := models.MatchCoach{
			CoachId: coachHome,
			Sanctions: models.Sanctions{
				Exclusions: []models.Exclusions{},
				YellowCard: false,
				RedCard:    false,
				BlueCard:   false,
				Report:     ""},
		}
		match_coaches_repository.CreateMatchCoach(match.AssociationId, matchCoach)
	}

	for _, coachAway := range startMatchRequest.CoachsAway {
		matchCoach := models.MatchCoach{
			CoachId: coachAway,
			Sanctions: models.Sanctions{
				Exclusions: []models.Exclusions{},
				YellowCard: false,
				RedCard:    false,
				BlueCard:   false,
				Report:     ""},
		}
		match_coaches_repository.CreateMatchCoach(match.AssociationId, matchCoach)
	}

	match.Referees = startMatchRequest.Referees
	match.Scorekeeper = startMatchRequest.Scorekeeper
	match.Timekeeper = startMatchRequest.Timekeeper

	return matches_repository.StartMatch(match, id)
}

func StartSecondHalf(id string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.FirstHalf {
		return false, errors.New("The match must be found in the first half")
	}

	return matches_repository.StartSecondHalf(id)
}

func EndMatch(id string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.SecondHalf {
		return false, errors.New("The match must be found in the second half")
	}

	return matches_repository.EndMatch(id)
}
