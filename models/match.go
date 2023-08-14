package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	//Status
	Created    = "created"
	Programmed = "programmed"
	FirstHalf  = "first_half"
	SecondHalf = "second_half"
	Finished   = "finished"

	//Sanctions
	Exclusion  = "exclusion"
	YellowCard = "yellow_card"
	RedCard    = "red_card"
	BlueCard   = "blue_card"
)

type Match struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date             time.Time          `bson:"date" json:"date"`
	TeamLocal        string             `bson:"team_local" json:"team_local"`
	TeamVisiting     string             `bson:"team_visiting" json:"team_visiting"`
	PlayersLocal     []string           `bson:"players_local" json:"players_local"`
	PlayersVisiting  []string           `bson:"players_visiting" json:"players_visiting"`
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
	AssociationId    string             `bson:"association_id" json:"association_id"`
	Status_Data
}

type TimeOut struct {
	Date time.Time `bson:"date" json:"date"`
	Half string    `bson:"half" json:"half"`
}

type MatchGoals struct {
	FirstHalf  int `bson:"first_half" json:"first_half"`
	SecondHalf int `bson:"second_half" json:"second_half"`
	Total      int `bson:"total" json:"total"`
}

type Sanction struct {
	Type     string    `bson:"type" json:"type"`
	Half     string    `bson:"half" json:"half"`
	Date     time.Time `bson:"date" json:"date"`
	Comments string    `bson:"comments" json:"comments"`
}

type Goal struct {
	Half string    `bson:"half" json:"half"`
	Date time.Time `bson:"date" json:"date"`
}

func (match *Match) SetCreatedDate() {
	match.CreatedDate = time.Now()
}

func (match *Match) SetModifiedDate() {
	match.ModifiedDate = time.Now()
}

func (match *Match) SetDisabled(disabled bool) {
	match.Disabled = disabled
}

func (match *Match) SetAssociationId(associationId string) {
	match.AssociationId = associationId
}

func generateMatch(teamA, teamB, phaseId string) Match {
	return Match{
		TeamLocal:    teamA,
		TeamVisiting: teamB,
		Status:       Created,
		PhaseId:      phaseId,
	}
}
