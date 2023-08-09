package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaguePhase struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams            []string           `bson:"teams,omitempty" json:"teams"`
	HomeAndAway      bool               `bson:"home_and_away,omitempty" json:"home_and_away"`
	ClassifiedNumber int                `bson:"classified_number,omitempty" json:"classified_number"`
	TeamsRanking     []TeamScore
}

type TeamScore struct {
	TeamId        string `bson:"team_id,omitempty" json:"team_id"`
	Points        int    `bson:"points,omitempty" json:"points"`
	Matches       int    `bson:"matches,omitempty" json:"matches"`
	Wins          int    `bson:"wins,omitempty" json:"wins"`
	Draws         int    `bson:"draws,omitempty" json:"draws"`
	Losses        int    `bson:"losses,omitempty" json:"losses"`
	GoalsScored   int    `bson:"goals_scored,omitempty" json:"goals_scored"`
	GoalsConceded int    `bson:"goals_conceded,omitempty" json:"goals_conceded"`
}
