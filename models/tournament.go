package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tournament struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Disabled      bool               `bson:"disabled" json:"disabled"`
	Status_Data   `bson:"status_data" json:"status_data"`
	AssociationId string `bson:"association_id" json:"association_id"`
}

func (tournament *Tournament) SetAssociationId(associationId string) {
	tournament.AssociationId = associationId
}

func (tournament *Tournament) SetCreatedDate() {
	tournament.CreatedDate = time.Now()
}

func (tournament *Tournament) SetModifiedDate() {
	tournament.ModifiedDate = time.Now()
}
