package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data       `bson:"personal_data" json:"personal_data"`
	AffiliateNumber     string    `bson:"affiliate_number" json:"affiliate_number"`
	TeamId              string    `bson:"team_id" json:"team_id"`
	ExpirationInsurance time.Time `bson:"expiration_insurance" json:"expiration_insurance"`
	AssociationId       string    `bson:"association_id" json:"association_id"`
	Status_Data         `bson:"status_data" json:"status_data"`
}

func (player *Player) SetCreatedDate() {
	player.CreatedDate = time.Now()
}

func (player *Player) SetModifiedDate() {
	player.ModifiedDate = time.Now()
}

func (player *Player) SetAssociationId(associationId string) {
	player.AssociationId = associationId
}

func (player *Player) SetId(id primitive.ObjectID) {
	player.Id = id
}

func (player *Player) SetAvatarURL(filename string) {
	player.Avatar = ImagesBaseURL + filename
}
