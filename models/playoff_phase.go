package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayoffPhase struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams                []TournamentTeamId `bson:"teams" json:"teams"`
	Config               PlayoffPhaseConfig `bson:"config" json:"config"`
	TournamentCategoryId string             `bson:"tournament_category_id" json:"tournament_category_id"`
	Status_Data          `bson:"status_data" json:"status_data"`
	AssociationId        string `bson:"association_id" json:"association_id"`
}

type PlayoffPhaseConfig struct {
	HomeAndAway      bool `bson:"home_and_away" json:"home_and_away"`
	SingleMatchFinal bool `bson:"single_match_final" json:"single_match_final"`
	RandomOrder      bool `bson:"random_order" json:"random_order"`
	ClassifiedNumber int  `bson:"classified_number" json:"classified_number"`
}

func (playoffPhase *PlayoffPhase) SetAssociationId(associationId string) {
	playoffPhase.AssociationId = associationId
}

func (playoffPhase *PlayoffPhase) SetCreatedDate() {
	playoffPhase.CreatedDate = time.Now()
}

func (playoffPhase *PlayoffPhase) SetModifiedDate() {
	playoffPhase.ModifiedDate = time.Now()
}

func (playoffPhase *PlayoffPhase) SetId(id primitive.ObjectID) {
	playoffPhase.Id = id
}
