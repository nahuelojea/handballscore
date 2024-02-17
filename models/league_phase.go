package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaguePhase struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams                []TournamentTeamId `bson:"teams" json:"teams"`
	TeamsRanking         []TeamScore        `bson:"teams_ranking" json:"teams_ranking"`
	Config               LeaguePhaseConfig  `bson:"config" json:"config"`
	TournamentCategoryId string             `bson:"tournament_category_id" json:"tournament_category_id"`
	Status_Data          `bson:"status_data" json:"status_data"`
	AssociationId        string `bson:"association_id" json:"association_id"`
}

type LeaguePhaseConfig struct {
	HomeAndAway      bool `bson:"home_and_away" json:"home_and_away"`
	ClassifiedNumber int  `bson:"classified_number" json:"classified_number"`
}

type TeamScore struct {
	TeamId        TournamentTeamId `bson:"team" json:"team"`
	Points        int              `bson:"points" json:"points"`
	Matches       int              `bson:"matches" json:"matches"`
	Wins          int              `bson:"wins" json:"wins"`
	Draws         int              `bson:"draws" json:"draws"`
	Losses        int              `bson:"losses" json:"losses"`
	GoalsScored   int              `bson:"goals_scored" json:"goals_scored"`
	GoalsConceded int              `bson:"goals_conceded" json:"goals_conceded"`
}

func (leaguePhase *LeaguePhase) SetAssociationId(associationId string) {
	leaguePhase.AssociationId = associationId
}

func (leaguePhase *LeaguePhase) SetCreatedDate() {
	leaguePhase.CreatedDate = time.Now()
}

func (leaguePhase *LeaguePhase) SetModifiedDate() {
	leaguePhase.ModifiedDate = time.Now()
}

func (leaguePhase *LeaguePhase) SetId(id primitive.ObjectID) {
	leaguePhase.Id = id
}

func (leaguePhase *LeaguePhase) InitializeTeamScores() {
	for _, teamId := range leaguePhase.Teams {
		teamScore := TeamScore{
			TeamId:        teamId,
			Points:        0,
			Matches:       0,
			Wins:          0,
			Draws:         0,
			Losses:        0,
			GoalsScored:   0,
			GoalsConceded: 0,
		}
		leaguePhase.TeamsRanking = append(leaguePhase.TeamsRanking, teamScore)
	}
}

func (leaguePhase LeaguePhase) GenerateLeaguePhaseWeeks() ([]LeaguePhaseWeek, [][]MatchRound) {
	var leaguePhaseWeeks []LeaguePhaseWeek

	rounds := calculateLeague(leaguePhase.Teams)

	var totalWeeks = len(rounds)

	if leaguePhase.Config.HomeAndAway {
		totalWeeks = totalWeeks * 2
	}

	for i := 0; i < totalWeeks; i++ {
		leaguePhaseWeek := LeaguePhaseWeek{
			Number:        i + 1,
			LeaguePhaseId: leaguePhase.Id.Hex(),
		}
		leaguePhaseWeeks = append(leaguePhaseWeeks, leaguePhaseWeek)
	}

	return leaguePhaseWeeks, rounds
}

func (leaguePhase LeaguePhase) GenerateMatches(rounds [][]MatchRound, leaguePhaseWeeks []LeaguePhaseWeek) []Match {
	var matches []Match
	var week = 1

	for i := 0; i < len(rounds); i++ {
		leaguePhaseWeek, _ := getLeaguePhaseWeekByNumber(leaguePhaseWeeks, week)

		for j := 0; j < len(rounds[i]); j++ {
			fmt.Printf("   %d-%d", 1+rounds[i][j].Home, 1+rounds[i][j].Away)
			matches = append(matches, generateLeagueMatch(leaguePhaseWeek.Id.Hex(),
				leaguePhase.Teams[rounds[i][j].Home],
				leaguePhase.Teams[rounds[i][j].Away]))
		}

		week++
		fmt.Println()
	}

	if leaguePhase.Config.HomeAndAway {
		for i := 0; i < len(rounds); i++ {
			leaguePhaseWeek, _ := getLeaguePhaseWeekByNumber(leaguePhaseWeeks, week)

			for j := 0; j < len(rounds[i]); j++ {
				fmt.Printf("   %d-%d", 1+rounds[i][j].Away, 1+rounds[i][j].Home)
				matches = append(matches, generateLeagueMatch(leaguePhaseWeek.Id.Hex(),
					leaguePhase.Teams[rounds[i][j].Away],
					leaguePhase.Teams[rounds[i][j].Home]))
			}

			fmt.Println()
		}
		week++
	}
	return matches
}

type MatchRound struct {
	Home int
	Away int
}

func calculateLeagueNumTeamsPair(teams []TournamentTeamId) [][]MatchRound {
	totalTeams := len(teams)
	totalRounds := totalTeams - 1
	totalMatchesByRound := totalTeams / 2

	rounds := make([][]MatchRound, totalRounds)
	for i := 0; i < totalRounds; i++ {
		rounds[i] = make([]MatchRound, totalMatchesByRound)
	}

	for i, k := 0, 0; i < totalRounds; i++ {
		for j := 0; j < totalMatchesByRound; j++ {
			rounds[i][j] = MatchRound{Home: k}

			k++

			if k == totalRounds {
				k = 0
			}
		}
	}

	for i := 0; i < totalRounds; i++ {
		if i%2 == 0 {
			rounds[i][0].Away = totalTeams - 1
		} else {
			rounds[i][0].Away = rounds[i][0].Home
			rounds[i][0].Home = totalTeams - 1
		}
	}

	teamHighest := totalTeams - 1
	teamOddHighest := teamHighest - 1

	for i, k := 0, teamOddHighest; i < totalRounds; i++ {
		for j := 1; j < totalMatchesByRound; j++ {
			rounds[i][j].Away = k

			k--

			if k == -1 {
				k = teamOddHighest
			}
		}
	}

	return rounds
}

func calculateLeagueNumTeamsOdd(teams []TournamentTeamId) [][]MatchRound {
	totalTeams := len(teams)
	totalRounds := totalTeams
	totalMatchesByRound := totalTeams / 2

	rounds := make([][]MatchRound, totalRounds)
	for i := 0; i < totalRounds; i++ {
		rounds[i] = make([]MatchRound, totalMatchesByRound)
	}

	for i, k := 0, 0; i < totalRounds; i++ {
		for j := -1; j < totalMatchesByRound; j++ {
			if j >= 0 {
				rounds[i][j] = MatchRound{Home: k}
			}

			k++

			if k == totalRounds {
				k = 0
			}
		}
	}

	teamHighest := totalTeams - 1

	for i, k := 0, teamHighest; i < totalRounds; i++ {
		for j := 0; j < totalMatchesByRound; j++ {
			rounds[i][j].Away = k

			k--

			if k == -1 {
				k = teamHighest
			}
		}
	}

	return rounds
}

func calculateLeague(teams []TournamentTeamId) [][]MatchRound {
	if len(teams)%2 == 0 {
		return calculateLeagueNumTeamsPair(teams)
	} else {
		return calculateLeagueNumTeamsOdd(teams)
	}
}

func getLeaguePhaseWeekByNumber(weeks []LeaguePhaseWeek, weekNumber int) (LeaguePhaseWeek, bool) {
	for _, week := range weeks {
		if week.Number == weekNumber {
			return week, true
		}
	}
	return LeaguePhaseWeek{}, false
}
