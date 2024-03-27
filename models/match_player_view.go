package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchPlayerView struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchId       string             `bson:"match_id" json:"match_id"`
	TeamId        TournamentTeamId   `bson:"team" json:"team"`
	PlayerId      string             `bson:"player_id" json:"player_id"`
	PlayerName    string             `bson:"player_name" json:"player_name"`
	PlayerSurname string             `bson:"player_surname" json:"player_surname"`
	PlayerAvatar  string             `bson:"player_avatar" json:"player_avatar"`
	Number        string             `bson:"number" json:"number"`
	Goals         `bson:"goals" json:"goals"`
	Sanctions     `bson:"sanctions" json:"sanctions"`
	AssociationId string `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
}

func (matchPlayerView *MatchPlayerView) SetCreatedDate() {
	matchPlayerView.CreatedDate = time.Now()
}

func (matchPlayerView *MatchPlayerView) SetModifiedDate() {
	matchPlayerView.ModifiedDate = time.Now()
}

func (matchPlayerView *MatchPlayerView) SetAssociationId(associationId string) {
	matchPlayerView.AssociationId = associationId
}

func (matchPlayerView *MatchPlayerView) SetId(id primitive.ObjectID) {
	matchPlayerView.Id = id
}
