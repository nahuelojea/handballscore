package matches_service

import (
	"errors"
	"fmt"
	"math"
	"time"

	dto "github.com/nahuelojea/handballscore/dto/matches"
	"github.com/nahuelojea/handballscore/handlers/end_match"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/match_coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/match_players_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
)

type GetMatchesOptions struct {
	TournamentCategoryId string
	LeaguePhaseWeekId  string
	PlayoffRoundKeyIds []string
	TeamId 		   string
	Variant 	   string
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
		TournamentCategoryId: filterOptions.TournamentCategoryId,
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

func GetMatchesTodayOrClosest(filterOptions GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
	filterOptions.Date = time.Now()
	matches, totalRecords, totalPages, err := GetMatchHeaders(filterOptions)
	if err != nil {
		return nil, 0, 0, err
	}

	if totalRecords > 0 {
		return matches, totalRecords, totalPages, nil
	}

	pastMatches, pastRecords, pastDays, err := getPastMatches(filterOptions)
	if err != nil {
		return nil, 0, 0, err
	}

	futureMatches, futureRecords, futureDays, err := getFutureMatches(filterOptions)
	if err != nil {
		return nil, 0, 0, err
	}

	if pastRecords == 0 && futureRecords == 0 {
		return nil, 0, 0, nil
	}

	if pastRecords > 0 && futureRecords > 0 {
		if pastDays == futureDays {
			return futureMatches, futureRecords, 1, nil
		} else if pastDays < futureDays {
			return pastMatches, pastRecords, 1, nil
		}
		return futureMatches, futureRecords, 1, nil
	}

	if pastRecords > 0 {
		return pastMatches, pastRecords, 1, nil
	}

	return futureMatches, futureRecords, 1, nil
}

func getPastMatches(filterOptions GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
	for i := 1; i <= 7; i++ {
		filterOptions.Date = time.Now().AddDate(0, 0, -i)
		matches, totalRecords, totalPages, err := GetMatchHeaders(filterOptions)
		if err != nil {
			return nil, 0, 0, err
		}
		if totalRecords > 0 {
			return matches, totalRecords, totalPages, nil
		}
	}
	return nil, 0, 0, nil
}

func getFutureMatches(filterOptions GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
	for i := 1; i <= 7; i++ {
		filterOptions.Date = time.Now().AddDate(0, 0, i)
		matches, totalRecords, totalPages, err := GetMatchHeaders(filterOptions)
		if err != nil {
			return nil, 0, 0, err
		}
		if totalRecords > 0 {
			return matches, totalRecords, totalPages, nil
		}
	}
	return nil, 0, 0, nil
}

func GetMatchesToday(filterOptions GetMatchesOptions, exactDate bool) ([]models.MatchHeaderView, int64, int, error) {
	if exactDate {
		filters := matches_repository.GetMatchesOptions{
			Date:          filterOptions.Date,
			AssociationId: filterOptions.AssociationId,
			Page:          filterOptions.Page,
			PageSize:      filterOptions.PageSize,
			SortField:     filterOptions.SortField,
			SortOrder:     filterOptions.SortOrder,
		}
		return matches_repository.GetMatchHeaders(filters)
	} else {
		return GetMatchesTodayOrClosest(filterOptions)
	}
}

func GetMatchesByJourney(filterOptions GetMatchesOptions) ([]dto.MatchResponse, int64, int, error) {
	filters := matches_repository.GetMatchesOptions{
		TournamentCategoryId: filterOptions.TournamentCategoryId,
		LeaguePhaseWeekId:  filterOptions.LeaguePhaseWeekId,
		PlayoffRoundKeyIds: filterOptions.PlayoffRoundKeyIds,
		AssociationId:      filterOptions.AssociationId,
		Page:               filterOptions.Page,
		PageSize:           filterOptions.PageSize,
		SortField:          filterOptions.SortField,
		SortOrder:          filterOptions.SortOrder,
	}

	matches, _, _, err := matches_repository.GetMatchHeaders(filters)
	if err != nil {
		return nil, 0, 0, errors.New("Error to get matches: " + err.Error())
	}

	var matchesJourney []dto.MatchResponse
	for _, match := range matches {
		homeMatchTeam := dto.MatchTeamResponse{
			TeamId: match.TeamHomeId.Hex(),
			Variant: match.TeamHomeVariant,
			Name:   match.TeamHomeName,
			Avatar: match.TeamHomeAvatar,
		}

		awayMatchTeam := dto.MatchTeamResponse{
			TeamId: match.TeamAwayId.Hex(),
			Variant: match.TeamAwayVariant,
			Name:   match.TeamAwayName,
			Avatar: match.TeamAwayAvatar,
		}

		matchJourney := dto.MatchResponse{
			MatchId:   match.Id.Hex(),
			Date:      match.Date,
			TeamHome:  homeMatchTeam,
			TeamAway:  awayMatchTeam,
			Referees:  match.Referees,
			Place:     match.Place,
			Status:    match.Status,
			GoalsHome: match.GoalsHome,
			GoalsAway: match.GoalsAway,
			PlayoffRound: match.PlayoffRound,
		}
		matchesJourney = append(matchesJourney, matchJourney)
	}

	totalRecords := int64(len(matchesJourney))

	totalPages := int(math.Ceil(float64(totalRecords) / float64(filterOptions.PageSize)))

	return matchesJourney, totalRecords, int(totalPages), nil
}

func GetMatchesByTeam(filterOptions GetMatchesOptions) ([]dto.MatchResponse, int64, int, error) {
	filters := matches_repository.GetMatchesOptions{
		TournamentCategoryId: filterOptions.TournamentCategoryId,
		TeamId:               filterOptions.TeamId,
		Variant:              filterOptions.Variant,
		AssociationId:      filterOptions.AssociationId,
		Page:               filterOptions.Page,
		PageSize:           filterOptions.PageSize,
		SortField:          filterOptions.SortField,
		SortOrder:          filterOptions.SortOrder,
	}

	matches, _, _, err := matches_repository.GetMatchHeaders(filters)
	if err != nil {
		return nil, 0, 0, errors.New("Error to get matches: " + err.Error())
	}

	var teamMatches []dto.MatchResponse
	for _, match := range matches {
		homeMatchTeam := dto.MatchTeamResponse{
			TeamId: match.TeamHomeId.Hex(),
			Variant: match.TeamHomeVariant,
			Name:   match.TeamHomeName,
			Avatar: match.TeamHomeAvatar,
		}

		awayMatchTeam := dto.MatchTeamResponse{
			TeamId: match.TeamAwayId.Hex(),
			Variant: match.TeamAwayVariant,
			Name:   match.TeamAwayName,
			Avatar: match.TeamAwayAvatar,
		}

		teamMatch := dto.MatchResponse{
			MatchId:   match.Id.Hex(),
			Date:      match.Date,
			TeamHome:  homeMatchTeam,
			TeamAway:  awayMatchTeam,
			Referees:  match.Referees,
			Place:     match.Place,
			Status:    match.Status,
			GoalsHome: match.GoalsHome,
			GoalsAway: match.GoalsAway,
			PlayoffRound: match.PlayoffRound,
		}
		teamMatches = append(teamMatches, teamMatch)
	}

	totalRecords := int64(len(teamMatches))

	totalPages := int(math.Ceil(float64(totalRecords) / float64(filterOptions.PageSize)))

	return teamMatches, totalRecords, int(totalPages), nil
}

func ProgramMatch(matchTime time.Time, place, id string) (bool, error) {
	if matchTime.Before(time.Now()) {
		return false, errors.New("The date cannot be earlier than the current date")
	}

	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	laodMatchPlayersAndCoachesFromLastMatch(match)

	return matches_repository.ProgramMatch(matchTime, place, id)
}

func laodMatchPlayersAndCoachesFromLastMatch(match models.Match) error {
	err := processTeamPlayersAndCoaches(match, match.TeamHome)
	if err != nil {
		return err
	}

	err = processTeamPlayersAndCoaches(match, match.TeamAway)
	if err != nil {
		return err
	}
	return nil
}

func processTeamPlayersAndCoaches(match models.Match, team models.TournamentTeamId) error {
	lastEndedMatch, isPresent, _ := matches_repository.GetLastEndedMatchByTeam(team, match.TournamentCategoryId)
	if !isPresent {
		return nil
	}

	if time.Now().Year() != lastEndedMatch.Date.Year() {
		return nil
	}

	if match.Status == models.Programmed || match.Status == models.Suspended {
		return nil
	}

	getPlayersOptions := match_players_repository.GetMatchPlayerOptions{
		MatchId:       lastEndedMatch.Id.Hex(),
		Team:          team,
		AssociationId: match.AssociationId,
	}

	players, _, _, _ := match_players_repository.GetMatchPlayers(getPlayersOptions)
	if len(players) != 0 {

		for _, player := range players {
			matchPlayer := models.MatchPlayer{
				MatchId:  match.Id.Hex(),
				PlayerId: player.PlayerId,
				Number:   player.Number,
				TeamId: models.TournamentTeamId{
					TeamId:  player.TeamId.TeamId,
					Variant: player.TeamId.Variant,
				},
				Goals: models.Goals{
					FirstHalf:  0,
					SecondHalf: 0,
				},
				Sanctions: models.Sanctions{
					Exclusions: []models.Exclusion{},
					YellowCard: false,
					RedCard:    false,
					BlueCard:   false,
					Report:     "",
				},
			}
			_, _, err := match_players_repository.CreateMatchPlayer(match.AssociationId, matchPlayer)
			if err != nil {
				fmt.Println("Error to create match player: %s", err.Error())
			}
		}
	}

	getCoachesOptions := match_coaches_repository.GetMatchCoachOptions{
		MatchId:       lastEndedMatch.Id.Hex(),
		Team:          team,
		AssociationId: match.AssociationId,
	}

	coaches, _, _, _ := match_coaches_repository.GetMatchCoaches(getCoachesOptions)
	if len(coaches) != 0 {
		for _, coach := range coaches {
			matchCoach := models.MatchCoach{
				MatchId: match.Id.Hex(),
				CoachId: coach.CoachId,
				TeamId: models.TournamentTeamId{
					TeamId:  coach.TeamId.TeamId,
					Variant: coach.TeamId.Variant,
				},
				Sanctions: models.Sanctions{
					Exclusions: []models.Exclusion{},
					YellowCard: false,
					RedCard:    false,
					BlueCard:   false,
					Report:     "",
				},
			}

			_, _, err := match_coaches_repository.CreateMatchCoach(match.AssociationId, matchCoach)
			if err != nil {
				fmt.Println("Error to create match coach: %s", err.Error())
			}
		}
	}

	return nil
}

func AssingReferees(id string, assignRefereesRequest dto.AssingRefereesRequest) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	if len(assignRefereesRequest.Referees) == 0 {
		return false, errors.New("The match must have at least one referee")
	}

	match.Referees = assignRefereesRequest.Referees

	return matches_repository.UpdateReferees(match, id)
}

func StartMatch(startMatchRequest dto.StartMatchRequest, id string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

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

	match.Status = models.Suspended
	err = end_match.EndMatchChainEvents(&match)
	if err != nil {
		return false, errors.New("Error to end match chain events: " + err.Error())
	}

	return true, nil
}

func EndMatch(id, comments string) (bool, error) {
	match, _, err := matches_repository.GetMatch(id)
	if err != nil {
		return false, errors.New("Error to get match: " + err.Error())
	}

	if match.Status != models.SecondHalf {
		return false, errors.New("The match must be found in the second half")
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
