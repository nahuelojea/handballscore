package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayoffPhase struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams  []string           `bson:"teams,omitempty" json:"teams"`
	Phases []LeaguePhase      `bson:"phases" json:"phases"`
}
