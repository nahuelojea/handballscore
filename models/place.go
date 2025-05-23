package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UbicationCoordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

type Place struct {
	Id            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Status_Data   `bson:"status_data" json:"status_data"`
	AssociationId string               `bson:"association_id" json:"association_id"`
	Name          string               `bson:"name" json:"name"`
	Ubication     UbicationCoordinates `bson:"ubication" json:"ubication"`
}

func (p *Place) SetCreatedDate() {
	p.Status_Data.CreatedDate = time.Now()
}

func (p *Place) SetModifiedDate() {
	p.Status_Data.ModifiedDate = time.Now()
}

func (p *Place) SetAssociationId(associationId string) {
	p.AssociationId = associationId
}

func (p *Place) SetId(id primitive.ObjectID) {
	p.Id = id
}
