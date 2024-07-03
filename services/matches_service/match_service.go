package matches_service

import (
	"errors"
	"math"
	"strings"
	"time"

	dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/handlers/end_match"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
	"github.com/nahuelojea/handballscore/services/teams_service"
)

type GetMatchesOptions struct {
	LeaguePhaseWeekId  string
	PlayoffRoundKeyIds []string
	Date               time.Time
	AssociationId      string
	Page               int
	PageSize           int
	SortField          string
	SortOrder          int
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

func GetMatchHeader(ID string) (models.MatchHeaderView, bool, error) {
	return matches_repository.GetMatchHeaderView(ID)
}

func GetMatches(filterOptions GetMatchesOptions) ([]models.Match, int64, int, error) {
	filters := matches_repository.GetMatchesOptions{
		LeaguePhaseWeekId:  filterOptions.LeaguePhaseWeekId,
		PlayoffRoundKeyIds: filterOptions.PlayoffRoundKeyIds,
		Date:               filterOptions.Date,
		AssociationId:      filterOptions.AssociationId,
		Page:               filterOptions.Page,
		PageSize:           filterOptions.PageSize,
		SortField:          filterOptions.SortField,
		SortOrder:          filterOptions.SortOrder,
	}
	return matches_repository.GetMatches(filters)
}

func GetMatchHeaders(filterOptions GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
	filters := matches_repository.GetMatchesOptions{
		Date:          filterOptions.Date,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return matches_repository.GetMatchHeaders(filters)
}

func GetMatchesByJourney(filterOptions GetMatchesOptions) ([]dto.MatchResponse, int64, int, error) {
	filters := matches_repository.GetMatchesOptions{
		LeaguePhaseWeekId:  filterOptions.LeaguePhaseWeekId,
		PlayoffRoundKeyIds: filterOptions.PlayoffRoundKeyIds,
		AssociationId:      filterOptions.AssociationId,
		Page:               filterOptions.Page,
		PageSize:           filterOptions.PageSize,
		SortField:          filterOptions.SortField,
		SortOrder:          filterOptions.SortOrder,
	}

	matches, _, _, err := matches_repository.GetMatches(filters)
	if err != nil {
		return nil, 0, 0, errors.New("Error to get matches: " + err.Error())
	}

	var matchesJourney []dto.MatchResponse
	for _, match := range matches {
		teamHome, _, err := teams_service.GetTeam(match.TeamHome.TeamId)
		if err != nil {
			return nil, 0, 0, errors.New("Error to get team: " + match.TeamHome.TeamId + " " + err.Error())
		}

		teamAway, _, err := teams_service.GetTeam(match.TeamAway.TeamId)
		if err != nil {
			return nil, 0, 0, errors.New("Error to get team: " + match.TeamAway.TeamId + " " + err.Error())
		}

		homeMatchTeam := dto.MatchTeamResponse{
			TeamId: teamHome.Id.Hex(),
			Name:   strings.TrimSpace(teamHome.Name + " " + match.TeamHome.Variant),
			Avatar: teamHome.Avatar,
		}

		awayMatchTeam := dto.MatchTeamResponse{
			TeamId: teamAway.Id.Hex(),
			Name:   strings.TrimSpace(teamAway.Name + " " + match.TeamAway.Variant),
			Avatar: teamAway.Avatar,
		}

		matchJourney := dto.MatchResponse{
			MatchId:           match.Id.Hex(),
			Date:              match.Date,
			TeamHome:          homeMatchTeam,
			TeamAway:          awayMatchTeam,
			Place:             match.Place,
			Status:            match.Status,
			AuthorizationCode: match.AuthorizationCode,
			GoalsHome:         match.GoalsHome.Total,
			GoalsAway:         match.GoalsAway.Total,
		}
		matchesJourney = append(matchesJourney, matchJourney)
	}

	totalRecords := int64(len(matchesJourney))

	totalPages := int(math.Ceil(float64(totalRecords) / float64(filterOptions.PageSize)))

	return matchesJourney, totalRecords, int(totalPages), nil
}

func ProgramMatch(matchTime time.Time, place, authorizationCode, id string) (bool, error) {
	if matchTime.Compare(time.Now()) < 1 {
		return false, errors.New("The date cannot be earlier than the current date")
	}
	return matches_repository.ProgramMatch(matchTime, place, authorizationCode, id)
}

func StartMatch(startMatchRequest dto.StartMatchRequest, id string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	for _, coachHome := range startMatchRequest.CoachsHome {
		matchCoach := models.MatchCoach{
			CoachId: coachHome,
			MatchId: match.Id.Hex(),
			TeamId:  match.TeamHome,
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
			MatchId: match.Id.Hex(),
			TeamId:  match.TeamAway,
			Sanctions: models.Sanctions{
				Exclusions: []models.Exclusion{},
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

func SuspendMatch(id, comments string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	_, err = matches_repository.SuspendMatch(id, comments)

	if err != nil {
		return false, errors.New("Error to suspend match: " + err.Error())
	}

	err = end_match.EndMatchChainEvents(&match)
	if err != nil {
		return false, errors.New("Error to end match chain events: " + err.Error())
	}

	return true, nil
}

func EndMatch(id, authorizationCode, comments string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.SecondHalf {
		return false, errors.New("The match must be found in the second half")
	}

	if match.AuthorizationCode != authorizationCode {
		return false, errors.New("The authorization code is incorrect")
	}

	_, err = matches_repository.EndMatch(id, comments)

	if err != nil {
		return false, errors.New("Error to end match: " + err.Error())
	}

	err = end_match.EndMatchChainEvents(&match)
	if err != nil {
		return false, errors.New("Error to end match chain events: " + err.Error())
	}

	return true, nil
}

func UpdateGoals(match models.Match, tournamentTeamId models.TournamentTeamId, add bool) (bool, error) {

	if match.TeamHome == tournamentTeamId || match.TeamAway == tournamentTeamId {
		if match.TeamHome == tournamentTeamId {
			if add {
				if match.Status == models.FirstHalf {
					match.GoalsHome.FirstHalf++
				} else {
					match.GoalsHome.SecondHalf++
				}
			} else {
				if match.Status == models.FirstHalf {
					if match.GoalsHome.FirstHalf > 0 {
						match.GoalsHome.FirstHalf--
					}
				} else {
					if match.GoalsHome.SecondHalf > 0 {
						match.GoalsHome.SecondHalf--
					}
				}
			}
		} else {
			if add {
				if match.Status == models.FirstHalf {
					match.GoalsAway.FirstHalf++
				} else {
					match.GoalsAway.SecondHalf++
				}
			} else {
				if match.Status == models.FirstHalf {
					if match.GoalsAway.FirstHalf > 0 {
						match.GoalsAway.FirstHalf--
					}
				} else {
					if match.GoalsAway.SecondHalf > 0 {
						match.GoalsAway.SecondHalf--
					}
				}
			}
		}
	} else {
		return false, errors.New("The team id does not match any of the two in the match")
	}

	match.GoalsHome.Total = match.GoalsHome.FirstHalf + match.GoalsHome.SecondHalf
	match.GoalsAway.Total = match.GoalsAway.FirstHalf + match.GoalsAway.SecondHalf

	return matches_repository.UpdateGoals(match, match.Id.Hex())
}

func UpdateTimeouts(id string, tournamentTeamId models.TournamentTeamId, add bool, time string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.FirstHalf && match.Status != models.SecondHalf {
		return false, errors.New("The match must be in play")
	}

	if match.Status == models.FirstHalf {
		time = "1º " + time
	} else {
		time = "2º " + time
	}

	if match.TeamHome == tournamentTeamId || match.TeamAway == tournamentTeamId {
		if match.TeamHome == tournamentTeamId {
			if add {
				if len(match.TimeoutsHome) == 3 {
					return false, errors.New("The team has three timeouts")
				}
				match.TimeoutsHome = append(match.TimeoutsHome, models.Timeout{Half: match.Status, Time: time})
			} else {
				if len(match.TimeoutsHome) > 0 {
					match.TimeoutsHome = match.TimeoutsHome[:len(match.TimeoutsHome)-1]
				}
			}
		} else {
			if add {
				if len(match.TimeoutsAway) == 3 {
					return false, errors.New("The team has three timeouts")
				}
				match.TimeoutsAway = append(match.TimeoutsAway, models.Timeout{Half: match.Status, Time: time})
			} else {
				if len(match.TimeoutsAway) > 0 {
					match.TimeoutsAway = match.TimeoutsAway[:len(match.TimeoutsAway)-1]
				}
			}
		}
	} else {
		return false, errors.New("The team id does not match any of the two in the match")
	}

	return matches_repository.UpdateTimeouts(match, id)
}

func GetPendingMatchesByLeaguePhaseId(leaguePhaseId string) ([]models.Match, error) {
	return matches_repository.GetPendingMatchesByLeaguePhaseId(leaguePhaseId)
}
