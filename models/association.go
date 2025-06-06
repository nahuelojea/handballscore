package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Association struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	DateOfFoundation time.Time          `bson:"date_of_foundation" json:"date_of_foundation"`
	Email            string             `bson:"email" json:"email"`
	Instagram        string             `bson:"instagram" json:"instagram"`
	News             []AssociationNew   `bson:"news" json:"news"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	PhoneNumber      string             `bson:"phone_number" json:"phone_number"`
}

type AssociationNew struct {
	Image string    `bson:"image" json:"image"`
	Date  time.Time `bson:"date" json:"date"`
}
