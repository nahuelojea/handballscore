package models

type Entity interface {
	SetCreatedDate()
	SetModifiedDate()
	SetDisabled(disabled bool)
	SetAssociationId(associationId string)
}
