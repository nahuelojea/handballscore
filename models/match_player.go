package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MatchPlayer struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchId   string             `bson:"match_id" json:"match_id"`
	PlayerId  string             `bson:"player_id" json:"player_id"`
	Sanctions []Sanction         `bson:"sanctions" json:"sanctions"`
	Goals     []Goal             `bson:"goals" json:"goals"`
}
