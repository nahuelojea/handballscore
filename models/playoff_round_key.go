package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayoffRoundKey struct {
	Id             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	KeyNumber      string              `bson:"key_number" json:"key_number"`
	Teams          [2]TournamentTeamId `bson:"tournament_team_id" json:"tournament_team_id"`
	TeamsRanking   [2]TeamScore        `bson:"teams_ranking" json:"teams_ranking"`
	PlayoffRoundId string              `bson:"playoff_round_id" json:"playoff_round_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
	AssociationId  string `bson:"association_id" json:"association_id"`
}

func (playoffRoundKey *PlayoffRoundKey) SetAssociationId(associationId string) {
	playoffRoundKey.AssociationId = associationId
}

func (playoffRoundKey *PlayoffRoundKey) SetCreatedDate() {
	playoffRoundKey.CreatedDate = time.Now()
}

func (playoffRoundKey *PlayoffRoundKey) SetModifiedDate() {
	playoffRoundKey.ModifiedDate = time.Now()
}

func (playoffRoundKey *PlayoffRoundKey) SetId(id primitive.ObjectID) {
	playoffRoundKey.Id = id
}