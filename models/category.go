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
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Gender        string             `bson:"gender" json:"gender"`
	AgeLimitFrom  int                `bson:"age_limit_from" json:"age_limit_from"`
	AgeLimitTo    int                `bson:"age_limit_to" json:"age_limit_to"`
	AssociationId string             `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
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
