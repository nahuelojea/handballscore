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
	CategoryId    string             `bson:"categoty_id" json:"categoty_id"`
	Teams         []string           `bson:"teams" json:"teams"`
	LeaguePhase   LeaguePhase        `bson:"league_phase" json:"league_phase"`
	PlayoffPhase  PlayoffPhase       `bson:"playoff_phase" json:"playoff_phase"`
	Status        string             `bson:"status" json:"status"`
	Champion      string             `bson:"champion" json:"champion"`
	TournamentId  string             `bson:"tournament_id" json:"tournament_id"`
	AssociationId string             `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
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

func (tournamentCategory *TournamentCategory) GenerateLeagueMatches() []Match {
	tournamentCategory.LeaguePhase.Id = primitive.NewObjectID()
	tournamentCategory.LeaguePhase.Teams = tournamentCategory.Teams
	return tournamentCategory.LeaguePhase.GenerateMatches()
}
