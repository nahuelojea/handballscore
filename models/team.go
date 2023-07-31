package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name,omitempty"`
	Address          string             `bson:"address" json:"address,omitempty"`
	PhoneNumber      string             `bson:"phone_number" json:"phone_number"`
	DateOfFoundation time.Time          `bson:"date_of_foundation" json:"date_of_foundation,omitempty"`
	Email            string             `bson:"email" json:"email"`
	Logo             string             `bson:"logo" json:"logo,omitempty"`
	AssociationId    string             `bson:"association_id" json:"association_id,omitempty"`
}
