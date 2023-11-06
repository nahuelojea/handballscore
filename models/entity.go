package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entity interface {
	SetCreatedDate()
	SetModifiedDate()
	SetAssociationId(associationId string)
	SetId(id primitive.ObjectID)
}
