package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tournament struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name,omitempty"`
	AssociationId string             `bson:"association_id" json:"association_id,omitempty"`
}
