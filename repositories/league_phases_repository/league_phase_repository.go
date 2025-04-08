package league_phases_repository

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	league_phase_weeks_repository "github.com/nahuelojea/handballscore/repositories/league_phases_weeks_repository"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
)

const (
	league_phase_collection = "league_phases"
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
	return repositories.Create(league_phase_collection, association_id, &leaguePhase)
}

func GetLeaguePhase(ID string) (models.LeaguePhase, bool, error) {
	var leaguePhase models.LeaguePhase
	_, err := repositories.GetById(league_phase_collection, ID, &leaguePhase)
	if err != nil {
		return models.LeaguePhase{}, false, err
	}

	return leaguePhase, true, nil
}

func GetLeaguePhases(filterOptions GetLeaguePhasesOptions) ([]models.LeaguePhase, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(league_phase_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.TournamentCategoryId != "" {
		filter["tournament_category_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.TournamentCategoryId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "tournament_category_id", Value: sortOrder},
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(sortFields)

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var leaguePhases []models.LeaguePhase
	for cur.Next(ctx) {
		var leaguePhase models.LeaguePhase
		if err := cur.Decode(&leaguePhase); err != nil {
			return nil, 0, 0, err
		}
		leaguePhases = append(leaguePhases, leaguePhase)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return leaguePhases, totalRecords, totalPages, nil
}

func UpdateTeamsRanking(leaguePhase models.LeaguePhase, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	if leaguePhase.TeamsRanking != nil {
		updateDataMap["teams_ranking"] = leaguePhase.TeamsRanking
	}

	return repositories.Update(league_phase_collection, updateDataMap, id)
}

func FinishPhase(id string, winner models.TournamentTeamId) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["finished"] = true
	updateDataMap["winner"] = winner

	return repositories.Update(league_phase_collection, updateDataMap, id)
}

func DeleteLeaguePhase(ID string) (bool, error) {
	return repositories.Delete(league_phase_collection, ID)
}

func ApplyOlympicTiebreaker(leaguePhase *models.LeaguePhase) {
	sort.SliceStable(leaguePhase.TeamsRanking, func(i, j int) bool {
		if leaguePhase.TeamsRanking[i].Points != leaguePhase.TeamsRanking[j].Points {
			return leaguePhase.TeamsRanking[i].Points > leaguePhase.TeamsRanking[j].Points
		}

		// Head-to-head result
		headToHeadResult := getHeadToHeadResult(leaguePhase.TeamsRanking[i].TeamId, leaguePhase.TeamsRanking[j].TeamId, *leaguePhase)
		if headToHeadResult != 0 {
			return headToHeadResult > 0
		}

		// Goal difference
		goalDifferenceA := leaguePhase.TeamsRanking[i].GoalsScored - leaguePhase.TeamsRanking[i].GoalsConceded
		goalDifferenceB := leaguePhase.TeamsRanking[j].GoalsScored - leaguePhase.TeamsRanking[j].GoalsConceded
		if goalDifferenceA != goalDifferenceB {
			return goalDifferenceA > goalDifferenceB
		}

		// Goals scored
		if leaguePhase.TeamsRanking[i].GoalsScored != leaguePhase.TeamsRanking[j].GoalsScored {
			return leaguePhase.TeamsRanking[i].GoalsScored > leaguePhase.TeamsRanking[j].GoalsScored
		}

		// Random draw
		return rand.Intn(2) == 0
	})
}

func getHeadToHeadResult(teamA, teamB models.TournamentTeamId, leaguePhase models.LeaguePhase) int {
	matches, _, _, err := matches_repository.GetMatchHeaders(matches_repository.GetMatchesOptions{
		Teams:                []models.TournamentTeamId{teamA, teamB},
		AssociationId:        leaguePhase.AssociationId,
		TournamentCategoryId: leaguePhase.TournamentCategoryId,
	})
	if err != nil {
		return 0
	}

	var teamAWins, teamBWins int
	for _, match := range matches {
		if match.Status == models.Finished {
			matchTeamHome := models.TournamentTeamId{TeamId: match.TeamHomeId.Hex(), Variant: match.TeamHomeVariant}
			matchTeamAway := models.TournamentTeamId{TeamId: match.TeamAwayId.Hex(), Variant: match.TeamAwayVariant}

			if matchTeamHome == teamA && matchTeamAway == teamB {
				if match.GoalsHome > match.GoalsAway {
					teamAWins++
				} else if match.GoalsHome < match.GoalsAway {
					teamBWins++
				}
			} else if matchTeamHome == teamB && matchTeamAway == teamA {
				if match.GoalsHome > match.GoalsAway {
					teamBWins++
				} else if match.GoalsHome < match.GoalsAway {
					teamAWins++
				}
			}
		}
	}

	if teamAWins > teamBWins {
		return 1
	} else if teamBWins > teamAWins {
		return -1
	}
	return 0
}

func RecalculateTeamsScores(leaguePhaseId string) error {
	leaguePhase, _, err := GetLeaguePhase(leaguePhaseId)
	if err != nil {
		return errors.New("League phase not found")
	}

	league_phases_weeks, _, _, err := league_phase_weeks_repository.GetLeaguePhaseWeeks(league_phase_weeks_repository.GetLeaguePhaseWeeksOptions{LeaguePhaseId: leaguePhaseId, AssociationId: leaguePhase.AssociationId})
	if err != nil {
		return errors.New("Error fetching league phase weeks")
	}

	if len(league_phases_weeks) > 0 {
		teamsScores := make([]models.TeamScore, 0)
		for _, league_phase_week := range league_phases_weeks {

			filterOptions := matches_repository.GetMatchesOptions{
				LeaguePhaseWeekId: league_phase_week.Id.Hex(),
				AssociationId:     league_phase_week.AssociationId,
				Page:              1,
				PageSize:          50,
				SortOrder:         1,
			}
			matches, _, _, err := matches_repository.GetMatches(filterOptions)
			if err != nil {
				return errors.New("Error fetching matches for league phase week")
			}

			for _, match := range matches {
				updateStandings(match, &teamsScores)
			}
		}
		leaguePhase.TeamsRanking = teamsScores

		_, err = UpdateTeamsRanking(leaguePhase, leaguePhase.Id.Hex())
		if err != nil {
			return errors.New("Error updating teams ranking in repository")
		}
	}

	return err
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
