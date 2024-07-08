package match_players_service

import (
	"errors"

	dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
	"github.com/nahuelojea/handballscore/services/matches_service"
)

func CreateMatchPlayer(association_id string, matchPlayerRequest dto.MatchPlayerRequest) (string, bool, error) {
	match, _, err := matches_repository.GetMatch(matchPlayerRequest.MatchId)
	if err != nil {
		return "", false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.Programmed &&
		match.Status != models.FirstHalf &&
		match.Status != models.SecondHalf {
		return "", false, errors.New("The player cannot be added in this match instance")
	}

	matchPlayer := models.MatchPlayer{
		MatchId:  matchPlayerRequest.MatchId,
		PlayerId: matchPlayerRequest.PlayerId,
		Number:   matchPlayerRequest.Number,
		TeamId: models.TournamentTeamId{
			TeamId:  matchPlayerRequest.Team.Id,
			Variant: matchPlayerRequest.Team.Variant,
		},
		Goals: models.Goals{
			FirstHalf:  0,
			SecondHalf: 0},
		Sanctions: models.Sanctions{
			Exclusions: []models.Exclusion{},
			YellowCard: false,
			RedCard:    false,
			BlueCard:   false,
			Report:     ""},
	}
	return match_players_repository.CreateMatchPlayer(association_id, matchPlayer)
}

func UpdateMatchPlayer(matchPlayer models.MatchPlayer, id string) (bool, error) {
	return match_players_repository.UpdateMatchPlayer(matchPlayer, id)
}

func DeleteMatchPlayer(id string) (bool, error) {
	return match_players_repository.DeleteMatchPlayer(id)
}

func GetMatchPlayer(id string) (models.MatchPlayer, bool, error) {
	return match_players_repository.GetMatchPlayer(id)
}

type GetMatchPlayerOptions struct {
	MatchId       string
	Team          models.TournamentTeamId
	PlayerId      string
	Number        string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchPlayers(filterOptions GetMatchPlayerOptions) ([]models.MatchPlayerView, int64, int, error) {
	filters := match_players_repository.GetMatchPlayerOptions{
		MatchId:       filterOptions.MatchId,
		Team:          filterOptions.Team,
		PlayerId:      filterOptions.PlayerId,
		Number:        filterOptions.Number,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}

	return match_players_repository.GetMatchPlayers(filters)
}

func UpdateGoal(id string, addGoal bool) (bool, error) {
	matchPlayer, match, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	if matchPlayer.RedCard {
		return false, errors.New("The player has red card")
	}

	if matchPlayer.BlueCard {
		return false, errors.New("The player has blue card")
	}

	matches_service.UpdateGoals(match, matchPlayer.TeamId, addGoal)

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

	matchPlayer.Goals.Total = matchPlayer.Goals.FirstHalf + matchPlayer.Goals.SecondHalf

	return match_players_repository.UpdateGoals(matchPlayer, match.Status)
}

func UpdateExclusions(id string, addExclusion bool, time string) (bool, error) {
	matchPlayer, match, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	if matchPlayer.RedCard {
		return false, errors.New("The player has red card")
	}

	if matchPlayer.BlueCard {
		return false, errors.New("The player has blue card")
	}

	if addExclusion {
		if len(matchPlayer.Sanctions.Exclusions) == 3 {
			return false, errors.New("The player has three exclusions")
		}
		if match.Status == models.FirstHalf {
			time = "1º " + time
		} else {
			time = "2º " + time
		}
		matchPlayer.Exclusions = append(matchPlayer.Exclusions, models.Exclusion{Time: time})
	} else {
		if len(matchPlayer.Exclusions) > 0 {
			matchPlayer.Exclusions = matchPlayer.Exclusions[:len(matchPlayer.Exclusions)-1]
		}
	}
	return match_players_repository.UpdateExclusions(matchPlayer)
}

func UpdateYellowCard(id string, addYellowCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	if matchPlayer.RedCard {
		return false, errors.New("The player has red card")
	}

	if matchPlayer.BlueCard {
		return false, errors.New("The player has blue card")
	}

	matchPlayer.YellowCard = addYellowCard

	return match_players_repository.UpdateYellowCard(matchPlayer)
}

func UpdateRedCard(id string, addRedCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
	if err != nil {
		return false, err
	}

	matchPlayer.RedCard = addRedCard

	return match_players_repository.UpdateRedCard(matchPlayer)
}

func UpdateNumber(id string, number int) (bool, error) {
	matchPlayer, _, err := match_players_repository.GetMatchPlayer(id)
	if err != nil {
		return false, errors.New("Error to get match player: " + err.Error())
	}

	matchPlayer.Number = number

	return match_players_repository.UpdateMatchPlayer(matchPlayer, id)
}

func UpdateBlueCard(id, report string, addBlueCard bool) (bool, error) {
	matchPlayer, _, err := getMatchPlayerAvailableToAction(id)
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

func getMatchPlayerAvailableToAction(id string) (models.MatchPlayer, models.Match, error) {
	matchPlayer, _, err := match_players_repository.GetMatchPlayer(id)
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

	return matchPlayer, match, nil
}
