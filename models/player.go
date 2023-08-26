package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	Personal_Data   `bson:"personal_data" json:"personal_data"`
	Gender          string `bson:"gender" json:"gender"`
	AffiliateNumber string `bson:"affiliate_number" json:"affiliate_number"`
	TeamId          string `bson:"team_id" json:"team_id"`
	AssociationId   string `bson:"association_id" json:"association_id"`
	Status_Data     `bson:"status_data" json:"status_data"`
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

func (player *Player) SetAssociationId(associationId string) {
	player.AssociationId = associationId
}
