package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coach struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	Status_Data
	TeamId        string `bson:"team_id" json:"team_id,omitempty"`
	AssociationId string `bson:"association_id" json:"association_id,omitempty"`
}

func (coach *Coach) SetCreatedDate() {
	coach.CreatedDate = time.Now()
}

func (coach *Coach) SetModifiedDate() {
	coach.ModifiedDate = time.Now()
}

func (coach *Coach) SetDisabled(disabled bool) {
	coach.Disabled = disabled
}

func (coach *Coach) GetAssociationId() string {
	return coach.AssociationId
}
