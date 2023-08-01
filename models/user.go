package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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
}
