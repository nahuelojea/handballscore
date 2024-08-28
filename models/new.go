package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type New struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssociationId string             `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
}

func (new *New) SetCreatedDate() {
	new.CreatedDate = time.Now()
}

func (new *New) SetModifiedDate() {
	new.ModifiedDate = time.Now()
}

func (new *New) SetAssociationId(associationId string) {
	new.AssociationId = associationId
}

func (new *New) SetId(id primitive.ObjectID) {
	new.Id = id
}
