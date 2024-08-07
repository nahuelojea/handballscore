package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AdminRole       string = "admin"
	AssociationRole string = "association"
	TeamRole        string = "team"
	RefereeRole     string = "referee"
	ViewerRole      string = "viewer"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"password,omitempty"`
	Role          string             `bson:"role" json:"role"`
	RoleId        string             `bson:"role_id" json:"role_id"`
	TeamId        string             `bson:"team_id" json:"team_id"`
	Personal_Data `bson:"personal_data" json:"personal_data"`
	Status_Data   `bson:"status_data" json:"status_data"`
	AssociationId string `bson:"association_id" json:"association_id"`
}

func (user *User) SetCreatedDate() {
	user.CreatedDate = time.Now()
}

func (user *User) SetModifiedDate() {
	user.ModifiedDate = time.Now()
}

func (user *User) SetDisabled(disabled bool) {
	user.Disabled = disabled
}

func (user *User) SetAssociationId(associationId string) {
	user.AssociationId = associationId
}

func (user *User) SetId(id primitive.ObjectID) {
	user.Id = id
}

func (user *User) SetAvatarURL(filename string) {
	user.Avatar = ImagesBaseURL + filename
}
