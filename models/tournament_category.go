package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Started = "started"
	Ended   = "ended"
)

type TournamentCategory struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	StartDate     time.Time          `bson:"start_date" json:"start_date"`
	EndDate       time.Time          `bson:"end_date" json:"end_date"`
	Status        string             `bson:"status" json:"status"`
	Teams         []TournamentTeamId `bson:"teams" json:"teams"`
	ChampionId    string             `bson:"champion_id" json:"champion_id"`
	TournamentId  string             `bson:"tournament_id" json:"tournament_id"`
	CategoryId    string             `bson:"category_id" json:"category_id"`
	AssociationId string             `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
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

type TournamentTeamId struct {
	TeamId  string `bson:"team_id" json:"team_id"`
	Variant string `bson:"variant" json:"variant"`
}

func (tournamentCategory *TournamentCategory) SetCreatedDate() {
	tournamentCategory.CreatedDate = time.Now()
}

func (tournamentCategory *TournamentCategory) SetModifiedDate() {
	tournamentCategory.ModifiedDate = time.Now()
}

func (tournamentCategory *TournamentCategory) SetAssociationId(associationId string) {
	tournamentCategory.AssociationId = associationId
}

func (tournamentCategory *TournamentCategory) SetId(id primitive.ObjectID) {
	tournamentCategory.Id = id
}
