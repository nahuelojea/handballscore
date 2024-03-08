package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchPlayer struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchId       string             `bson:"match_id" json:"match_id"`
	TeamId        TournamentTeamId   `bson:"team_id" json:"team_id"`
	PlayerId      string             `bson:"player_id" json:"player_id"`
	Number        string             `bson:"number" json:"number"`
	Goals         `bson:"goals" json:"goals"`
	Sanctions     `bson:"sanctions" json:"sanctions"`
	AssociationId string `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
}

type Goals struct {
	FirstHalf  int `bson:"first_half" json:"first_half"`
	SecondHalf int `bson:"second_half" json:"second_half"`
}

func (matchPlayer *MatchPlayer) SetCreatedDate() {
	matchPlayer.CreatedDate = time.Now()
}

func (matchPlayer *MatchPlayer) SetModifiedDate() {
	matchPlayer.ModifiedDate = time.Now()
}

func (matchPlayer *MatchPlayer) SetAssociationId(associationId string) {
	matchPlayer.AssociationId = associationId
}

func (matchPlayer *MatchPlayer) SetId(id primitive.ObjectID) {
	matchPlayer.Id = id
}
