package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Surname     string             `bson:"surname" json:"surname,omitempty"`
	DateOfBirth time.Time          `bson:"date_of_birth" json:"date_of_birth,omitempty"`
	Dni         string             `bson:"dni" json:"dni,omitempty"`
	Avatar      string             `bson:"avatar" json:"avatar,omitempty"`
}
