package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaguePhaseWeek struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number        int                `bson:"number" json:"number"`
	LeaguePhaseId string             `bson:"league_phase_id" json:"league_phase_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
	AssociationId string `bson:"association_id" json:"association_id"`
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetAssociationId(associationId string) {
	leaguePhaseWeek.AssociationId = associationId
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetCreatedDate() {
	leaguePhaseWeek.CreatedDate = time.Now()
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetModifiedDate() {
	leaguePhaseWeek.ModifiedDate = time.Now()
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetId(id primitive.ObjectID) {
	leaguePhaseWeek.Id = id
}
