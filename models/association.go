package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Association struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name,omitempty"`
	DateOfFoundation time.Time          `bson:"date_of_foundation" json:"date_of_foundation,omitempty"`
	Email            string             `bson:"email" json:"email"`
	Logo             string             `bson:"logo" json:"logo,omitempty"`
}
