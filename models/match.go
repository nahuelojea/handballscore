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
)

type Match struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date               time.Time          `bson:"date" json:"date"`
	TeamHome           string             `bson:"team_home" json:"team_home"`
	TeamAway           string             `bson:"team_away" json:"team_away"`
	Referees           []string           `bson:"referees" json:"referees"`
	Place              string             `bson:"place" json:"place"`
	Scorekeeper        string             `bson:"scorekeeper" json:"scorekeeper"`
	Timekeeper         string             `bson:"timekeeper" json:"timekeeper"`
	Comments           string             `bson:"comments" json:"comments"`
	Status             string             `bson:"status" json:"status"`
	GoalsHome          MatchGoals         `bson:"goals_home" json:"goals_home"`
	GoalsAway          MatchGoals         `bson:"goals_away" json:"goals_away"`
	TimeoutsHome       []TimeOuts         `bson:"timeouts_home" json:"timeouts_home"`
	TimeoutsAway       []TimeOuts         `bson:"timeouts_away" json:"timeouts_away"`
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
	Exclusions int    `bson:"exclusions" json:"exclusions"`
	YellowCard bool   `bson:"yellow_card" json:"yellow_card"`
	RedCard    bool   `bson:"red_card" json:"red_card"`
	BlueCard   bool   `bson:"blue_card" json:"blue_card"`
	Report     string `bson:"report" json:"report"`
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
		TeamHome:          teamA,
		TeamAway:          teamB,
		Status:            Created,
		LeaguePhaseWeekId: leaguePhaseWeekId,
	}
}

func generatePlayoffMatch(playoffPhaseWeekId string, teamA, teamB string) Match {
	return Match{
		TeamHome:           teamA,
		TeamAway:           teamB,
		Status:             Created,
		PlayoffPhaseWeekId: playoffPhaseWeekId,
	}
}
