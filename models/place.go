package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UbicationCoordinates defines the latitude and longitude of a place.
type UbicationCoordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

// Place defines the structure for a place where matches can be played.
type Place struct {
	Id              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Status_Data     Status_Data          `bson:"status_data" json:"status_data"`
	AssociationId   string               `bson:"association_id" json:"association_id"`
	Name            string               `bson:"name" json:"name"`
	Ubication       UbicationCoordinates `bson:"ubication" json:"ubication"`
}

// SetCreatedDate sets the created date of the place.
func (p *Place) SetCreatedDate() {
	p.Status_Data.CreatedDate = time.Now()
}

// SetModifiedDate sets the modified date of the place.
func (p *Place) SetModifiedDate() {
	p.Status_Data.ModifiedDate = time.Now()
}

// SetAssociationId sets the association ID of the place.
func (p *Place) SetAssociationId(associationId string) {
	p.AssociationId = associationId
}

// SetId sets the ID of the place.
func (p *Place) SetId(id primitive.ObjectID) {
	p.Id = id
}
