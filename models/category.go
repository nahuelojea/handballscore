package models

import (
	"time"

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
	Status_Data
}

func (category *Category) SetCreatedDate() {
	category.CreatedDate = time.Now()
}

func (category *Category) SetModifiedDate() {
	category.ModifiedDate = time.Now()
}

func (category *Category) SetDisabled(disabled bool) {
	category.Disabled = disabled
}

func (category *Category) SetAssociationId(associationId string) {
	category.AssociationId = associationId
}
