package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Admin  Role = "admin"
	Viewer Role = "viewer"
)

type Role string

type User struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"password,omitempty"`
	Role          Role               `bson:"role" json:"role"`
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

func (user *User) SetAvatarURL(filename string) {
	user.Avatar = ImagesBaseURL + filename
}
