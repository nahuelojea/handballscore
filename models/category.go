package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Male   = "M"
	Female = "F"
)

type Category struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name,omitempty"`
	Gender        string             `bson:"gender" json:"gender,omitempty"`
	AgeLimitFrom  int                `bson:"age_limit_from" json:"age_limit_from,omitempty"`
	AgeLimitTo    int                `bson:"age_limit_to" json:"age_limit_to,omitempty"`
	AssociationId string             `bson:"association_id" json:"association_id,omitempty"`
}
