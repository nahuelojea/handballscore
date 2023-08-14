package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Started = "started"
	Ended   = "ended"
)

type Tournament struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name,omitempty"`
	CategoryId    string             `bson:"categoty_id" json:"categoty_id,omitempty"`
	Teams         []string           `bson:"teams" json:"teams,omitempty"`
	LeaguePhase   LeaguePhase        `bson:"league_phase" json:"league_phase,omitempty"`
	PlayoffPhase  PlayoffPhase       `bson:"playoff_phase" json:"playoff_phase,omitempty"`
	Status        string             `bson:"status" json:"status,omitempty"`
	Champion      string             `bson:"champion" json:"champion,omitempty"`
	AssociationId string             `bson:"association_id" json:"association_id,omitempty"`
	Status_Data
}

func (tournament *Tournament) SetCreatedDate() {
	tournament.CreatedDate = time.Now()
}

func (tournament *Tournament) SetModifiedDate() {
	tournament.ModifiedDate = time.Now()
}

func (tournament *Tournament) SetDisabled(disabled bool) {
	tournament.Disabled = disabled
}

func (tournament *Tournament) SetAssociationId(associationId string) {
	tournament.AssociationId = associationId
}

func (tournament *Tournament) GenerateLeagueMatches() []Match {
	tournament.LeaguePhase.Id = primitive.NewObjectID()
	tournament.LeaguePhase.Teams = tournament.Teams
	return tournament.LeaguePhase.GenerateMatches()
}
