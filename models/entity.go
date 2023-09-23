package models

type Entity interface {
	SetCreatedDate()
	SetModifiedDate()
	SetAssociationId(associationId string)
}
