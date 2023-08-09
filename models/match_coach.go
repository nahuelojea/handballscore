package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MatchCoach struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchId   string             `bson:"match_id" json:"match_id"`
	CoachId   string             `bson:"coach_id" json:"coach_id"`
	Sanctions []Sanction         `bson:"sanctions" json:"sanctions"`
}
