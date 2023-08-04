package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	Gender          string `bson:"gender" json:"gender,omitempty"`
	AffiliateNumber string `bson:"affiliate_number" json:"affiliate_number,omitempty"`
	TeamId          string `bson:"team_id" json:"team_id,omitempty"`
	AssociationId   string `bson:"association_id" json:"association_id,omitempty"`
	Status_Data
}

func (player *Player) SetCreatedDate() {
	player.CreatedDate = time.Now()
}

func (player *Player) SetModifiedDate() {
	player.ModifiedDate = time.Now()
}

func (player *Player) SetDisabled(disabled bool) {
	player.Disabled = disabled
}

func (player *Player) GetAssociationId() string {
	return player.AssociationId
}
