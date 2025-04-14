package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Address          string             `bson:"address" json:"address"`
	PhoneNumber      string             `bson:"phone_number" json:"phone_number"`
	DateOfFoundation time.Time          `bson:"date_of_foundation" json:"date_of_foundation"`
	Email            string             `bson:"email" json:"email"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	Initials         string             `bson:"initials" json:"initials"`
	MainColor        string             `bson:"main_color" json:"main_color"`
	SecondaryColor   string             `bson:"secondary_color" json:"secondary_color"`
	AssociationId    string             `bson:"association_id" json:"association_id"`
	Status_Data      `bson:"status_data" json:"status_data"`
}

func (team *Team) SetCreatedDate() {
	team.CreatedDate = time.Now()
}

func (team *Team) SetModifiedDate() {
	team.ModifiedDate = time.Now()
}

func (team *Team) SetAssociationId(associationId string) {
	team.AssociationId = associationId
}

func (team *Team) SetId(id primitive.ObjectID) {
	team.Id = id
}

func (team *Team) SetAvatarURL(filename string) {
	team.Avatar = ImagesBaseURL + filename
}
