package league_phases_service

import (
	"errors"
	"fmt"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phases_repository"
	"github.com/nahuelojea/handballscore/services/league_phase_weeks_service"
	"github.com/nahuelojea/handballscore/services/matches_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLeaguePhasesOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func CreateLeaguePhase(association_id string, leaguePhase models.LeaguePhase) (string, bool, error) {
	return league_phases_repository.CreateLeaguePhase(association_id, leaguePhase)
}

func GetLeaguePhase(ID string) (models.LeaguePhase, bool, error) {
	return league_phases_repository.GetLeaguePhase(ID)
}

func GetLeaguePhases(filterOptions GetLeaguePhasesOptions) ([]models.LeaguePhase, int64, int, error) {
	filters := league_phases_repository.GetLeaguePhasesOptions{
		TournamentCategoryId: filterOptions.TournamentCategoryId,
		AssociationId:        filterOptions.AssociationId,
		Page:                 filterOptions.Page,
		PageSize:             filterOptions.PageSize,
		SortField:            filterOptions.SortField,
		SortOrder:            filterOptions.SortOrder,
	}
	return league_phases_repository.GetLeaguePhases(filters)
}

func DeleteLeaguePhase(ID string) (bool, error) {
	return league_phases_repository.DeleteLeaguePhase(ID)
}

func CreateTournamentLeaguePhase(tournamentCategory models.TournamentCategory, leaguePhaseConfig models.LeaguePhaseConfig) (string, bool, error) {
	var leaguePhase models.LeaguePhase

	leaguePhase.TournamentCategoryId = tournamentCategory.Id.Hex()
	leaguePhase.Config = leaguePhaseConfig

	leaguePhase.Teams = tournamentCategory.Teams

	leaguePhase.InitializeTeamScores()

	leaguePhaseIdStr, _, err := CreateLeaguePhase(tournamentCategory.AssociationId, leaguePhase)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase: %s", err.Error()))
	}

	leaguePhaseId, err := primitive.ObjectIDFromHex(leaguePhaseIdStr)
	if err != nil {
		return "", false, err
	}

	leaguePhase.Id = leaguePhaseId
	leaguePhaseWeeks, rounds := leaguePhase.GenerateLeaguePhaseWeeks()

	_, _, err = league_phase_weeks_service.CreateLeaguePhaseWeeks(tournamentCategory.AssociationId, leaguePhaseWeeks)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase weeks: %s", err.Error()))
	}

	filterOptions := league_phase_weeks_service.GetLeaguePhaseWeeksOptions{
		AssociationId: tournamentCategory.AssociationId,
		LeaguePhaseId: leaguePhaseId.Hex(),
	}

	leaguePhaseWeeks, _, _, err = league_phase_weeks_service.GetLeaguePhaseWeeks(filterOptions)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to get league phase weeks: %s", err.Error()))
	}

	matches := leaguePhase.GenerateMatches(tournamentCategory.Id.Hex(), rounds, leaguePhaseWeeks)

	_, _, err = matches_service.CreateMatches(tournamentCategory.AssociationId, matches)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create league phase matches: %s", err.Error()))
	}

	return tournamentCategory.Id.Hex(), true, nil
}

func RecalculateTeamsScores(leaguePhaseId string) error {
	leaguePhase, _, err := league_phases_repository.GetLeaguePhase(leaguePhaseId)
	if err != nil {
		return errors.New("League phase not found")
	}

	league_phases_weeks, _, _, err := league_phase_weeks_service.GetLeaguePhaseWeeks(league_phase_weeks_service.GetLeaguePhaseWeeksOptions{LeaguePhaseId: leaguePhaseId, AssociationId: leaguePhase.AssociationId})
	if err != nil {
		return errors.New("Error fetching league phase weeks")
	}

	if len(league_phases_weeks) > 0 {
		teamsScores := make([]models.TeamScore, 0)
		for _, league_phase_week := range league_phases_weeks {

			filterOptions := matches_service.GetMatchesOptions{
				LeaguePhaseWeekId: league_phase_week.Id.Hex(),
				AssociationId:     league_phase_week.AssociationId,
				Page:              1,
				PageSize:          50,
				SortOrder:         1,
			}
			matches, _, _, err := matches_service.GetMatches(filterOptions)
			if err != nil {
				return errors.New("Error fetching matches for league phase week")
			}

			for _, match := range matches {
				updateStandings(match, &teamsScores)
			}
		}
		leaguePhase.TeamsRanking = teamsScores

		_, err = league_phases_repository.UpdateTeamsRanking(leaguePhase, leaguePhase.Id.Hex())
		if err != nil {
			return errors.New("Error updating teams ranking in repository")
		}
	}

	return err
}

func ApplyOlympicTiebreaker (leaguePhase *models.LeaguePhase) {
	league_phases_repository.ApplyOlympicTiebreaker(leaguePhase)
}

func updateStandings(match models.Match, teamsScores *[]models.TeamScore) {
	homeTeamScore := findTeamInStandings(match.TeamHome, teamsScores)
	awayTeamScore := findTeamInStandings(match.TeamAway, teamsScores)

	if match.Status == models.Finished {
		updateTeamScores(match, homeTeamScore, awayTeamScore)

		updateTeamInSlice(match.TeamHome, homeTeamScore, teamsScores)
		updateTeamInSlice(match.TeamAway, awayTeamScore, teamsScores)
	}
}

func updateTeamInSlice(teamId models.TournamentTeamId, updatedScore *models.TeamScore, standings *[]models.TeamScore) {
	for i := range *standings {
		if (*standings)[i].TeamId == teamId {
			(*standings)[i] = *updatedScore
			break
		}
	}
}

func updateTeamScores(match models.Match, homeTeam *models.TeamScore, awayTeam *models.TeamScore) {
	switch {
	case match.GoalsHome.Total > match.GoalsAway.Total:
		if !(match.GoalsHome.Total == 9 && match.GoalsAway.Total == 0) {
			awayTeam.Points++
		}
		homeTeam.Points += 3
		homeTeam.Wins++
		awayTeam.Losses++
	case match.GoalsHome.Total < match.GoalsAway.Total:
		if !(match.GoalsAway.Total == 9 && match.GoalsHome.Total == 0) {
			homeTeam.Points++
		}
		awayTeam.Points += 3
		awayTeam.Wins++
		homeTeam.Losses++
	default:
		homeTeam.Points += 2
		awayTeam.Points += 2
		homeTeam.Draws++
		awayTeam.Draws++
	}

	homeTeam.GoalsScored += match.GoalsHome.Total
	homeTeam.GoalsConceded += match.GoalsAway.Total
	awayTeam.GoalsScored += match.GoalsAway.Total
	awayTeam.GoalsConceded += match.GoalsHome.Total

	homeTeam.Matches++
	awayTeam.Matches++
}

func findTeamInStandings(teamId models.TournamentTeamId, standings *[]models.TeamScore) *models.TeamScore {
	for i := range *standings {
		if (*standings)[i].TeamId == teamId {
			return &(*standings)[i]
		}
	}

	newTeamScore := models.TeamScore{
		TeamId: teamId,
	}
	*standings = append(*standings, newTeamScore)

	return &(*standings)[len(*standings)-1]
}
