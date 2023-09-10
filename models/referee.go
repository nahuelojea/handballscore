package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Referee struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data `bson:"personal_data" json:"personal_data"`
	AssociationId string `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
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
