package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tournament struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name,omitempty"`
	CategoryId    string             `bson:"categoty_id" json:"categoty_id,omitempty"`
	Teams         []string           `bson:"teams" json:"teams,omitempty"`
	LeaguePhase   LeaguePhase        `bson:"league_phase" json:"league_phase,omitempty"`
	PlayoffPhase  PlayoffPhase       `bson:"playoff_phase" json:"playoff_phase,omitempty"`
	Champion      string             `bson:"champion" json:"champion,omitempty"`
	AssociationId string             `bson:"association_id" json:"association_id,omitempty"`
}
