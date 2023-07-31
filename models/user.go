package models

type User struct {
	Person
	Email         string `bson:"email" json:"email"`
	Password      string `bson:"password" json:"password,omitempty"`
	AssociationId string `bson:"association_id" json:"association_id,omitempty"`
}
