package match_coaches_service

import (
	"errors"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

func CreateMatchCoach(association_id string, matchCoach models.MatchCoach) (string, bool, error) {
	return match_coaches_repository.CreateMatchCoach(association_id, matchCoach)
}

func DeleteMatchCoach(id string) (bool, error) {
	return match_coaches_repository.DeleteMatchCoach(id)
}

func GetMatchCoach(associationId, id string) (models.MatchCoach, bool, error) {
	return match_coaches_repository.GetMatchCoach(id)
}

type GetMatchCoachOptions struct {
	MatchId       string
	TeamId        string
	CoachId       string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchCoachs(filterOptions GetMatchCoachOptions) ([]models.MatchCoach, int64, error) {
	filters := match_coaches_repository.GetMatchCoachOptions{
		MatchId:       filterOptions.MatchId,
		TeamId:        filterOptions.TeamId,
		CoachId:       filterOptions.CoachId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}

	return match_coaches_repository.GetMatchCoaches(filters)
}

func UpdateExclusions(matchCoachId string, addExclusion bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(matchCoachId)
	if err != nil {
		return false, err
	}

	if addExclusion {
		if matchCoach.Exclusions == 2 {
			return false, errors.New("The coach has two exclusions")
		} else {
			matchCoach.Exclusions++
		}
	} else {
		if matchCoach.Exclusions != 0 {
			matchCoach.Exclusions--
		}
	}

	return match_coaches_repository.UpdateExclusions(matchCoach)
}

func UpdateYellowCard(matchCoachId string, addYellowCard bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(matchCoachId)
	if err != nil {
		return false, err
	}

	matchCoach.YellowCard = addYellowCard

	return match_coaches_repository.UpdateYellowCard(matchCoach)
}

func UpdateRedCard(matchCoachId string, addRedCard bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(matchCoachId)
	if err != nil {
		return false, err
	}

	matchCoach.RedCard = addRedCard

	return match_coaches_repository.UpdateRedCard(matchCoach)
}

func UpdateBlueCard(matchCoachId, report string, addBlueCard bool) (bool, error) {
	matchCoach, _, err := getMatchCoachAvailableToAction(matchCoachId)
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

func getMatchCoachAvailableToAction(matchCoachId string) (models.MatchCoach, models.Match, error) {
	matchCoach, _, err := match_coaches_repository.GetMatchCoach(matchCoachId)
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
