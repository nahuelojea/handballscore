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
	TeamLocal        MatchTeam          `bson:"team_local" json:"team_local"`
	TeamVisiting     MatchTeam          `bson:"team_visiting" json:"team_visiting"`
	PlayersLocal     []MatchPlayer      `bson:"players_local" json:"players_local"`
	CoachsLocal      []MatchCoach       `bson:"coachs_local" json:"coachs_local"`
	PlayersVisiting  []MatchPlayer      `bson:"players_visiting" json:"players_visiting"`
	CoachsVisiting   []MatchCoach       `bson:"coachs_visiting" json:"coachs_visiting"`
	Referees         []MatchReferee     `bson:"referees" json:"referees"`
	Place            string             `bson:"place" json:"place"`
	Scorekeeper      string             `bson:"scorekeeper" json:"scorekeeper"`
	Timekeeper       string             `bson:"timekeeper" json:"timekeeper"`
	Comments         string             `bson:"comments" json:"comments"`
	Status           string             `bson:"status" json:"status"`
	GoalsLocal       MatchGoals         `bson:"goals_local" json:"goals_local"`
	GoalsVisiting    MatchGoals         `bson:"goals_visiting" json:"goals_visiting"`
	TimeoutsLocal    []TimeOuts         `bson:"timeouts_local" json:"timeouts_local"`
	TimeoutsVisiting []TimeOuts         `bson:"timeouts_visiting" json:"timeouts_visiting"`
	PhaseId          string             `bson:"phase_id" json:"phase_id"`
	AssociationId    string             `bson:"association_id" json:"association_id"`
	Status_Data      `bson:"status_data" json:"status_data"`
}

type MatchPlayer struct {
	Id              string `bson:"_id" json:"id"`
	Name            string `bson:"name" json:"name"`
	AffiliateNumber string `bson:"affiliate_number" json:"affiliate_number"`
	Avatar          string `bson:"avatar" json:"avatar"`
	Number          string `bson:"number" json:"number"`
	Goals           `bson:"goals" json:"goals"`
	Sanctions       `bson:"sanctions" json:"sanctions"`
}

type MatchCoach struct {
	Id        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Avatar    string `bson:"avatar" json:"avatar"`
	Sanctions `bson:"sanctions" json:"sanctions"`
}

type MatchReferee struct {
	Id     string `bson:"_id" json:"id"`
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
}

type MatchTeam struct {
	Id     string `bson:"_id" json:"id"`
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
}

type Sanctions struct {
	Exclusions  int    `bson:"exclusions" json:"exclusions"`
	YellowCards int    `bson:"yellow_cards" json:"yellow_cards"`
	RedCard     bool   `bson:"red_card" json:"red_card"`
	BlueCard    bool   `bson:"blue_card" json:"blue_card"`
	Comments    string `bson:"comments" json:"comments"`
}

type Goals struct {
	FirstHalf  int `bson:"first_half" json:"first_half"`
	SecondHalf int `bson:"second_half" json:"second_half"`
}

type TimeOuts struct {
	Half string `bson:"half" json:"half"`
	Time string `bson:"time" json:"time"`
}

type MatchGoals struct {
	FirstHalf  int `bson:"first_half" json:"first_half"`
	SecondHalf int `bson:"second_half" json:"second_half"`
	Total      int `bson:"total" json:"total"`
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

func generateMatch(phaseId string, teamA, teamB MatchTeam) Match {
	return Match{
		TeamLocal:    teamA,
		TeamVisiting: teamB,
		Status:       Created,
		PhaseId:      phaseId,
	}
}
