package models

type Referee struct {
	Person
	AssociationId string `bson:"association_id" json:"association_id,omitempty"`
}
