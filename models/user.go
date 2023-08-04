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
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data
	Email         string `bson:"email" json:"email"`
	Password      string `bson:"password" json:"password,omitempty"`
	Role          Role   `bson:"role" json:"role,omitempty"`
	AssociationId string `bson:"association_id" json:"association_id,omitempty"`
	Status_Data
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

func (user *User) GetAssociationId() string {
	return user.AssociationId
}
