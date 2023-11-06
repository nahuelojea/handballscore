package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayoffPhase struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HomeAndAway bool               `bson:"home_and_away" json:"home_and_away"`
	RandomOrder bool               `bson:"random_order" json:"random_order"`
	Teams       []string           `bson:"teams" json:"teams"`
}
