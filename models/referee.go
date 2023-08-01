package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Referee struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	Status_Data
	AssociationId string `bson:"association_id" json:"association_id,omitempty"`
}
