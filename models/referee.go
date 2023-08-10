package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Referee struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	AssociationId string `bson:"association_id" json:"association_id,omitempty"`
	Status_Data
}

func (referee *Referee) SetCreatedDate() {
	referee.CreatedDate = time.Now()
}

func (referee *Referee) SetModifiedDate() {
	referee.ModifiedDate = time.Now()
}

func (referee *Referee) SetDisabled(disabled bool) {
	referee.Disabled = disabled
}

func (referee *Referee) SetAssociationId(associationId string) {
	referee.AssociationId = associationId
}
