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
	TeamHome           string             `bson:"team_home" json:"team_home"`
	TeamAway           string             `bson:"team_away" json:"team_away"`
	PlayersHome        []MatchPlayer      `bson:"players_home" json:"players_home"`
	CoachsHome         []MatchCoach       `bson:"coachs_home" json:"coachs_home"`
	PlayersAway        []MatchPlayer      `bson:"players_away" json:"players_away"`
	CoachsAway         []MatchCoach       `bson:"coachs_away" json:"coachs_away"`
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

type MatchCoach struct {
	CoachId   string `bson:"coach_id" json:"coach_id"`
	Sanctions `bson:"sanctions" json:"sanctions"`
}

type MatchPlayer struct {
	PlayerId  string `bson:"player_id" json:"player_id"`
	Number    string `bson:"number" json:"number"`
	Goals     `bson:"goals" json:"goals"`
	Sanctions `bson:"sanctions" json:"sanctions"`
}

type Goals struct {
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
