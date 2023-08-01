package models

import (
	"time"
)

type Personal_Data struct {
	Name        string    `bson:"name" json:"name,omitempty"`
	Surname     string    `bson:"surname" json:"surname,omitempty"`
	DateOfBirth time.Time `bson:"date_of_birth" json:"date_of_birth,omitempty"`
	Dni         string    `bson:"dni" json:"dni,omitempty"`
	PhoneNumber string    `bson:"phone_number" json:"phone_number"`
	Avatar      string    `bson:"avatar" json:"avatar,omitempty"`
}
