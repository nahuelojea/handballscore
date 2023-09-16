package models

import (
	"time"
)

type Personal_Data struct {
	Name        string    `bson:"name" json:"name"`
	Surname     string    `bson:"surname" json:"surname"`
	DateOfBirth time.Time `bson:"date_of_birth" json:"date_of_birth"`
	Dni         string    `bson:"dni" json:"dni"`
	Gender      string    `bson:"gender" json:"gender"`
	PhoneNumber string    `bson:"phone_number" json:"phone_number"`
	Avatar      string    `bson:"avatar" json:"avatar"`
}
