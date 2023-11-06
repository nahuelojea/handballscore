package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	//Sanctions
	Exclusion  = "exclusion"
	YellowCard = "yellow_card"
	RedCard    = "red_card"
	BlueCard   = "blue_card"

	//Status
	Created    = "created"
	Programmed = "programmed"
	FirstHalf  = "first_half"
	SecondHalf = "second_half"
	Finished   = "finished"
)

type Match struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date               time.Time          `bson:"date" json:"date"`
	TeamLocal          string             `bson:"team_local" json:"team_local"`
	TeamVisiting       string             `bson:"team_visiting" json:"team_visiting"`
	CoachsLocal        []string           `bson:"coachs_local" json:"coachs_local"`
	PlayersVisiting    []string           `bson:"players_visiting" json:"players_visiting"`
	CoachsVisiting     []string           `bson:"coachs_visiting" json:"coachs_visiting"`
	Referees           []string           `bson:"referees" json:"referees"`
	Place              string             `bson:"place" json:"place"`
	Scorekeeper        string             `bson:"scorekeeper" json:"scorekeeper"`
	Timekeeper         string             `bson:"timekeeper" json:"timekeeper"`
	Comments           string             `bson:"comments" json:"comments"`
	Status             string             `bson:"status" json:"status"`
	GoalsLocal         MatchGoals         `bson:"goals_local" json:"goals_local"`
	GoalsVisiting      MatchGoals         `bson:"goals_visiting" json:"goals_visiting"`
	TimeoutsLocal      []TimeOuts         `bson:"timeouts_local" json:"timeouts_local"`
	TimeoutsVisiting   []TimeOuts         `bson:"timeouts_visiting" json:"timeouts_visiting"`
	LeaguePhaseWeekId  string             `bson:"league_phase_week_id" json:"league_phase_week_id"`
	PlayoffPhaseWeekId string             `bson:"playoff_phase_week_id" json:"playoff_phase_week_id"`
	AssociationId      string             `bson:"association_id" json:"association_id"`
	Status_Data        `bson:"status_data" json:"status_data"`
}

type TimeOuts struct {
	Half string `bson:"half" json:"half"`
	Time string `bson:"time" json:"time"`
}

type MatchGoals struct {
	FirstHalf  int `bson:"first_half" json:"first_half"`
	SecondHalf int `bson:"second_half" json:"second_half"`
}

type Sanctions struct {
	Exclusions  int    `bson:"exclusions" json:"exclusions"`
	YellowCards int    `bson:"yellow_cards" json:"yellow_cards"`
	RedCard     bool   `bson:"red_card" json:"red_card"`
	BlueCard    bool   `bson:"blue_card" json:"blue_card"`
	Comments    string `bson:"comments" json:"comments"`
}

func (match *Match) SetCreatedDate() {
	match.CreatedDate = time.Now()
}

func (match *Match) SetModifiedDate() {
	match.ModifiedDate = time.Now()
}

func (match *Match) SetAssociationId(associationId string) {
	match.AssociationId = associationId
}

func (match *Match) SetId(id primitive.ObjectID) {
	match.Id = id
}

func generateLeagueMatch(leaguePhaseWeekId string, teamA, teamB string) Match {
	return Match{
		TeamLocal:         teamA,
		TeamVisiting:      teamB,
		Status:            Created,
		LeaguePhaseWeekId: leaguePhaseWeekId,
	}
}

func generatePlayoffMatch(playoffPhaseWeekId string, teamA, teamB string) Match {
	return Match{
		TeamLocal:          teamA,
		TeamVisiting:       teamB,
		Status:             Created,
		PlayoffPhaseWeekId: playoffPhaseWeekId,
	}
}
