package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Coach struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	Status_Data
	TeamId string `bson:"team_id" json:"team_id,omitempty"`
}
