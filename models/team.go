package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name,omitempty"`
	Address          string             `bson:"address" json:"address,omitempty"`
	PhoneNumber      string             `bson:"phone_number" json:"phone_number"`
	DateOfFoundation time.Time          `bson:"date_of_foundation" json:"date_of_foundation,omitempty"`
	Email            string             `bson:"email" json:"email"`
	Avatar           string             `bson:"avatar" json:"avatar,omitempty"`
	AssociationId    string             `bson:"association_id" json:"association_id,omitempty"`
	Status_Data
}

func (team *Team) SetCreatedDate() {
	team.CreatedDate = time.Now()
}

func (team *Team) SetModifiedDate() {
	team.ModifiedDate = time.Now()
}

func (team *Team) SetDisabled(disabled bool) {
	team.Disabled = disabled
}

func (team *Team) GetAssociationId() string {
	return team.AssociationId
}
