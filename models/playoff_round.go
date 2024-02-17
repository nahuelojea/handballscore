package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ThirtyTwoFinals = "thirty_two_finals"
	SixteenFinals   = "sixteen_finals"
	EightFinals     = "eight_finals"
	QuarterFinals   = "quarter_finals"
	SemiFinal       = "semi_final"
	Final           = "final"
)

type PlayoffRound struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Round          string             `bson:"round" json:"round"`
	PlayoffPhaseId string             `bson:"playoff_phase_id" json:"playoff_phase_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
	AssociationId  string `bson:"association_id" json:"association_id"`
}
