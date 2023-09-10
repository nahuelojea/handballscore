package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayoffPhase struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HomeAndAway bool               `bson:"home_and_away" json:"home_and_away"`
	Teams       []string           `bson:"teams" json:"teams"`
	Phases      []LeaguePhase      `bson:"phases" json:"phases"`
}
