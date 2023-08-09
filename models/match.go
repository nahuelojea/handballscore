package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date             time.Time          `bson:"date" json:"date"`
	TeamLocal        string             `bson:"team_local" json:"team_local"`
	TeamVisiting     string             `bson:"team_visiting" json:"team_visiting"`
	Referees         []string           `bson:"referees" json:"referees"`
	Place            string             `bson:"place" json:"place"`
	Scorekeeper      string             `bson:"scorekeeper" json:"scorekeeper"`
	Timekeeper       string             `bson:"timekeeper" json:"timekeeper"`
	Comments         string             `bson:"comments" json:"comments"`
	Status           string             `bson:"status" json:"status"`
	GoalsLocal       MatchGoals         `bson:"goals_local" json:"goals_local"`
	GoalsVisiting    MatchGoals         `bson:"goals_visiting" json:"goals_visiting"`
	TimeoutsLocal    []TimeOut          `bson:"timeouts_local" json:"timeouts_local"`
	TimeoutsVisiting []TimeOut          `bson:"timeouts_visiting" json:"timeouts_visiting"`
	PhaseId          string             `bson:"phase_id" json:"phase_id"`
}

type TimeOut struct {
	Date time.Time `bson:"date" json:"date"`
	Time string    `bson:"time" json:"time"`
}

type MatchGoals struct {
	FirstTime  int `bson:"first_time" json:"first_time"`
	SecondTime int `bson:"second_time" json:"second_time"`
	Total      int `bson:"total" json:"total"`
}

type Sanction struct {
	Type     string    `bson:"type" json:"type"`
	Time     string    `bson:"time" json:"time"`
	Date     time.Time `bson:"date" json:"date"`
	Comments string    `bson:"comments" json:"comments"`
}

type Goal struct {
	Time string    `bson:"time" json:"time"`
	Date time.Time `bson:"date" json:"date"`
}
