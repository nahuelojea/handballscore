package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaguePhase struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams            []MatchTeam        `bson:"teams,omitempty" json:"teams"`
	HomeAndAway      bool               `bson:"home_and_away,omitempty" json:"home_and_away"`
	ClassifiedNumber int                `bson:"classified_number,omitempty" json:"classified_number"`
	TeamsRanking     []TeamScore        `bson:"teams_ranking,omitempty" json:"teams_ranking"`
}

type TeamScore struct {
	Team          MatchTeam `bson:"team,omitempty" json:"team"`
	Points        int       `bson:"points,omitempty" json:"points"`
	Matches       int       `bson:"matches,omitempty" json:"matches"`
	Wins          int       `bson:"wins,omitempty" json:"wins"`
	Draws         int       `bson:"draws,omitempty" json:"draws"`
	Losses        int       `bson:"losses,omitempty" json:"losses"`
	GoalsScored   int       `bson:"goals_scored,omitempty" json:"goals_scored"`
	GoalsConceded int       `bson:"goals_conceded,omitempty" json:"goals_conceded"`
}

func (leaguePhase *LeaguePhase) GenerateMatches() []Match {
	var matches []Match

	if leaguePhase.HomeAndAway {
		for i, teamA := range leaguePhase.Teams {
			for j := i + 1; j < len(leaguePhase.Teams); j++ {
				teamB := leaguePhase.Teams[j]

				matches = append(matches, generateMatch(leaguePhase.Id.Hex(), teamA, teamB))
				matches = append(matches, generateMatch(leaguePhase.Id.Hex(), teamB, teamA))
			}
		}
	} else {
		totalTeams := len(leaguePhase.Teams)

		for i := 0; i < totalTeams-1; i++ {
			for j := i + 1; j < totalTeams; j++ {
				var local, visiting MatchTeam

				if (i+j)%2 == 0 {
					local, visiting = leaguePhase.Teams[i], leaguePhase.Teams[j]
				} else {
					local, visiting = leaguePhase.Teams[j], leaguePhase.Teams[i]
				}

				matches = append(matches, generateMatch(leaguePhase.Id.Hex(), local, visiting))
			}
		}
	}

	return matches
}
