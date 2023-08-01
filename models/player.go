package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	Status_Data
	Gender          string `bson:"gender" json:"gender,omitempty"`
	AffiliateNumber string `bson:"affiliate_number" json:"affiliate_number,omitempty"`
	TeamId          string `bson:"team_id" json:"team_id,omitempty"`
}
