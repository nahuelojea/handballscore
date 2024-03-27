package match_coaches_service

import (
	"errors"

	dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

func CreateMatchCoaches(startMatchRequest dto.StartMatchRequest, match models.Match) {
	for _, coachHome := range startMatchRequest.CoachsHome {
		matchCoach := models.MatchCoach{
			CoachId: coachHome,
			Sanctions: models.Sanctions{
				Exclusions: []models.Exclusion{},
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
				Exclusions: []models.Exclusion{},
				YellowCard: false,
				RedCard:    false,
				BlueCard:   false,
				Report:     ""},
		}
		match_coaches_repository.CreateMatchCoach(match.AssociationId, matchCoach)
	}
}

func CreateMatchCoach(association_id string, matchCoach models.MatchCoach) (string, bool, error) {
	return match_coaches_repository.CreateMatchCoach(association_id, matchCoach)
}

func DeleteMatchCoach(id string) (bool, error) {
	return match_coaches_repository.DeleteMatchCoach(id)
}

func GetMatchCoach(id string) (models.MatchCoach, bool, error) {
	return match_coaches_repository.GetMatchCoach(id)
}

type GetMatchCoachOptions struct {
	MatchId       string
	Team          models.TournamentTeamId
	CoachId       string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchCoaches(filterOptions GetMatchCoachOptions) ([]models.MatchCoachView, int64, int, error) {
	filters := match_coaches_repository.GetMatchCoachOptions{
		MatchId:       filterOptions.MatchId,
		Team:          filterOptions.Team,
		CoachId:       filterOptions.CoachId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}

	return match_coaches_repository.GetMatchCoaches(filters)
}

func UpdateExclusions(id string, addExclusion bool, time string) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(id)
	if err != nil {
		return false, err
	}

	if addExclusion {
		if len(matchCoach.Exclusions) == 2 {
			return false, errors.New("The coach has two exclusions")
		}
		matchCoach.Exclusions = append(matchCoach.Exclusions, models.Exclusion{Time: time})
	} else {
		if len(matchCoach.Exclusions) > 0 {
			matchCoach.Exclusions = matchCoach.Exclusions[:len(matchCoach.Exclusions)-1]
		}
	}
	return match_coaches_repository.UpdateExclusions(matchCoach)
}

func UpdateYellowCard(id string, addYellowCard bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchCoach.YellowCard = addYellowCard

	return match_coaches_repository.UpdateYellowCard(matchCoach)
}

func UpdateRedCard(id string, addRedCard bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchCoach.RedCard = addRedCard

	return match_coaches_repository.UpdateRedCard(matchCoach)
}

func UpdateBlueCard(id, report string, addBlueCard bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchCoach.BlueCard = addBlueCard
	if addBlueCard {
		matchCoach.Report = report
	} else {
		matchCoach.Report = ""
	}

	return match_coaches_repository.UpdateBlueCard(matchCoach)
}

func getMatchCoachAvailableToAction(id string) (models.MatchCoach, models.Match, error) {
	matchCoach, _, err := match_coaches_repository.GetMatchCoach(id)
	if err != nil {
		return models.MatchCoach{}, models.Match{}, errors.New("Error to get match coach: " + err.Error())
	}

	match, _, err := matches_repository.GetMatch(matchCoach.MatchId)
	if err != nil {
		return models.MatchCoach{}, models.Match{}, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.FirstHalf && match.Status != models.SecondHalf {
		return models.MatchCoach{}, models.Match{}, errors.New("The match must be in progress")
	}

	if matchCoach.RedCard {
		return models.MatchCoach{}, models.Match{}, errors.New("The coach has red card")
	}

	if matchCoach.BlueCard {
		return models.MatchCoach{}, models.Match{}, errors.New("The coach has blue card")
	}

	return matchCoach, match, nil
}
