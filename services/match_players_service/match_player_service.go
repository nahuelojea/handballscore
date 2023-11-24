package match_players_service

import (
	"errors"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

func CreateMatchPlayers(association_id string, matchPlayers []models.MatchPlayer) ([]string, bool, error) {
	return match_players_repository.CreateMatchPlayers(association_id, matchPlayers)
}

func CreateMatchPlayer(association_id string, matchPlayer models.MatchPlayer) (string, bool, error) {
	return match_players_repository.CreateMatchPlayer(association_id, matchPlayer)
}

func UpdateMatchPlayer(matchPlayer models.MatchPlayer, id string) (bool, error) {
	return match_players_repository.UpdateMatchPlayer(matchPlayer, id)
}

func DeleteMatchPlayer(id string) (bool, error) {
	return match_players_repository.DeleteMatchPlayer(id)
}

func GetMatchPlayer(associationId, id string) (models.MatchPlayer, bool, error) {
	return match_players_repository.GetMatchPlayer(id)
}

type GetMatchPlayerOptions struct {
	MatchId       string
	TeamId        string
	PlayerId      string
	Number        string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchPlayers(filterOptions GetMatchPlayerOptions) ([]models.MatchPlayer, int64, error) {
	filters := match_players_repository.GetMatchPlayerOptions{
		MatchId:       filterOptions.MatchId,
		TeamId:        filterOptions.TeamId,
		PlayerId:      filterOptions.PlayerId,
		Number:        filterOptions.Number,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}

	return match_players_repository.GetMatchPlayers(filters)
}

func UpdateGoal(matchPlayerId string, addGoal bool) (bool, error) {
	matchPlayer, match, err := getMatchPlayerAvailableToAction(matchPlayerId)
	if err != nil {
		return false, err
	}

	if match.Status == models.FirstHalf {
		if addGoal {
			matchPlayer.Goals.FirstHalf++
		} else {
			if matchPlayer.Goals.FirstHalf != 0 {
				matchPlayer.Goals.FirstHalf--
			}
		}
	} else {
		if addGoal {
			matchPlayer.Goals.SecondHalf++
		} else {
			if matchPlayer.Goals.SecondHalf != 0 {
				matchPlayer.Goals.SecondHalf--
			}
		}
	}

	return match_players_repository.UpdateGoals(matchPlayer, match.Status)
}

func UpdateExclusions(matchPlayerId string, addExclusion bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(matchPlayerId)
	if err != nil {
		return false, err
	}

	if addExclusion {
		if matchPlayer.Exclusions == 2 {
			return false, errors.New("The player has two exclusions")
		} else {
			matchPlayer.Exclusions++
		}
	} else {
		if matchPlayer.Exclusions != 0 {
			matchPlayer.Exclusions--
		}
	}

	return match_players_repository.UpdateExclusions(matchPlayer)
}

func UpdateYellowCard(matchPlayerId string, addYellowCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(matchPlayerId)
	if err != nil {
		return false, err
	}

	matchPlayer.YellowCard = addYellowCard

	return match_players_repository.UpdateYellowCard(matchPlayer)
}

func UpdateRedCard(matchPlayerId string, addRedCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(matchPlayerId)
	if err != nil {
		return false, err
	}

	matchPlayer.RedCard = addRedCard

	return match_players_repository.UpdateRedCard(matchPlayer)
}

func UpdateBlueCard(matchPlayerId, report string, addBlueCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(matchPlayerId)
	if err != nil {
		return false, err
	}

	matchPlayer.BlueCard = addBlueCard
	if addBlueCard {
		matchPlayer.Report = report
	} else {
		matchPlayer.Report = ""
	}

	return match_players_repository.UpdateBlueCard(matchPlayer)
}

func getMatchPlayerAvailableToAction(matchPlayerId string) (models.MatchPlayer, models.Match, error) {
	matchPlayer, _, err := match_players_repository.GetMatchPlayer(matchPlayerId)
	if err != nil {
		return models.MatchPlayer{}, models.Match{}, errors.New("Error to get match player: " + err.Error())
	}

	match, _, err := matches_repository.GetMatch(matchPlayer.MatchId)
	if err != nil {
		return models.MatchPlayer{}, models.Match{}, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.FirstHalf && match.Status != models.SecondHalf {
		return models.MatchPlayer{}, models.Match{}, errors.New("The match must be in progress")
	}

	if matchPlayer.RedCard {
		return models.MatchPlayer{}, models.Match{}, errors.New("The player has red card")
	}

	if matchPlayer.BlueCard {
		return models.MatchPlayer{}, models.Match{}, errors.New("The player has blue card")
	}

	return matchPlayer, match, nil
}
