package models

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPlace_SetCreatedDate(t *testing.T) {
	place := &Place{}
	place.SetCreatedDate()
	if place.Status_Data.CreatedDate.IsZero() {
		t.Errorf("Expected CreatedDate to be set, but it was zero")
	}
	// Check if the date is recent (e.g., within the last few seconds)
	if time.Since(place.Status_Data.CreatedDate) > 5*time.Second {
		t.Errorf("CreatedDate is not recent: %s", place.Status_Data.CreatedDate)
	}
}

func TestPlace_SetModifiedDate(t *testing.T) {
	place := &Place{}
	place.SetModifiedDate()
	if place.Status_Data.ModifiedDate.IsZero() {
		t.Errorf("Expected ModifiedDate to be set, but it was zero")
	}
	if time.Since(place.Status_Data.ModifiedDate) > 5*time.Second {
		t.Errorf("ModifiedDate is not recent: %s", place.Status_Data.ModifiedDate)
	}
}

func TestPlace_SetAssociationId(t *testing.T) {
	place := &Place{}
	testAssociationId := "assoc123"
	place.SetAssociationId(testAssociationId)
	if place.AssociationId != testAssociationId {
		t.Errorf("Expected AssociationId to be %s, but got %s", testAssociationId, place.AssociationId)
	}
}

func TestPlace_SetId(t *testing.T) {
	place := &Place{}
	testId := primitive.NewObjectID()
	place.SetId(testId)
	if place.Id != testId {
		t.Errorf("Expected Id to be %s, but got %s", testId.Hex(), place.Id.Hex())
	}
}

func TestUbicationCoordinates_BSONTags(t *testing.T) {
	uc := UbicationCoordinates{}
	ucType := reflect.TypeOf(uc)

	latitudeField, _ := ucType.FieldByName("Latitude")
	if tag := latitudeField.Tag.Get("bson"); tag != "latitude" {
		t.Errorf("Expected bson tag 'latitude' for Latitude field, got '%s'", tag)
	}
	if tag := latitudeField.Tag.Get("json"); tag != "latitude" {
		t.Errorf("Expected json tag 'latitude' for Latitude field, got '%s'", tag)
	}

	longitudeField, _ := ucType.FieldByName("Longitude")
	if tag := longitudeField.Tag.Get("bson"); tag != "longitude" {
		t.Errorf("Expected bson tag 'longitude' for Longitude field, got '%s'", tag)
	}
	if tag := longitudeField.Tag.Get("json"); tag != "longitude" {
		t.Errorf("Expected json tag 'longitude' for Longitude field, got '%s'", tag)
	}
}

func TestPlace_BSONTags(t *testing.T) {
	p := Place{}
	pType := reflect.TypeOf(p)

	idField, _ := pType.FieldByName("Id")
	if tag := idField.Tag.Get("bson"); tag != "_id,omitempty" {
		t.Errorf("Expected bson tag '_id,omitempty' for Id field, got '%s'", tag)
	}
	if tag := idField.Tag.Get("json"); tag != "id" {
		t.Errorf("Expected json tag 'id' for Id field, got '%s'", tag)
	}

	statusDataField, _ := pType.FieldByName("Status_Data")
	if tag := statusDataField.Tag.Get("bson"); tag != "status_data" {
		t.Errorf("Expected bson tag 'status_data' for Status_Data field, got '%s'", tag)
	}
	if tag := statusDataField.Tag.Get("json"); tag != "status_data" {
		t.Errorf("Expected json tag 'status_data' for Status_Data field, got '%s'", tag)
	}

	associationIdField, _ := pType.FieldByName("AssociationId")
	if tag := associationIdField.Tag.Get("bson"); tag != "association_id" {
		t.Errorf("Expected bson tag 'association_id' for AssociationId field, got '%s'", tag)
	}
	if tag := associationIdField.Tag.Get("json"); tag != "association_id" {
		t.Errorf("Expected json tag 'association_id' for AssociationId field, got '%s'", tag)
	}

	nameField, _ := pType.FieldByName("Name")
	if tag := nameField.Tag.Get("bson"); tag != "name" {
		t.Errorf("Expected bson tag 'name' for Name field, got '%s'", tag)
	}
	if tag := nameField.Tag.Get("json"); tag != "name" {
		t.Errorf("Expected json tag 'name' for Name field, got '%s'", tag)
	}

	ubicationField, _ := pType.FieldByName("Ubication")
	if tag := ubicationField.Tag.Get("bson"); tag != "ubication" {
		t.Errorf("Expected bson tag 'ubication' for Ubication field, got '%s'", tag)
	}
	if tag := ubicationField.Tag.Get("json"); tag != "ubication" {
		t.Errorf("Expected json tag 'ubication' for Ubication field, got '%s'", tag)
	}
}

// Test BSON marshaling for Place (selected fields)
func TestPlace_BSONMarshalling(t *testing.T) {
	id := primitive.NewObjectID()
	assocId := "assocTest123"
	placeName := "Test Place"
	lat := 12.34
	lon := 56.78

	place := Place{
		Id:            id,
		AssociationId: assocId,
		Name:          placeName,
		Ubication: UbicationCoordinates{
			Latitude:  lat,
			Longitude: lon,
		},
	}
	place.SetCreatedDate()

	bytes, err := bson.Marshal(place)
	if err != nil {
		t.Fatalf("BSON marshaling failed: %v", err)
	}

	var unmarshaledPlace Place
	err = bson.Unmarshal(bytes, &unmarshaledPlace)
	if err != nil {
		t.Fatalf("BSON unmarshaling failed: %v", err)
	}

	if unmarshaledPlace.Id != id {
		t.Errorf("Expected Id %s, got %s", id.Hex(), unmarshaledPlace.Id.Hex())
	}
	if unmarshaledPlace.AssociationId != assocId {
		t.Errorf("Expected AssociationId %s, got %s", assocId, unmarshaledPlace.AssociationId)
	}
	if unmarshaledPlace.Name != placeName {
		t.Errorf("Expected Name %s, got %s", placeName, unmarshaledPlace.Name)
	}
	if unmarshaledPlace.Ubication.Latitude != lat {
		t.Errorf("Expected Latitude %f, got %f", lat, unmarshaledPlace.Ubication.Latitude)
	}
	if unmarshaledPlace.Ubication.Longitude != lon {
		t.Errorf("Expected Longitude %f, got %f", lon, unmarshaledPlace.Ubication.Longitude)
	}
	if unmarshaledPlace.Status_Data.CreatedDate.IsZero() {
		t.Error("Expected CreatedDate to be set after unmarshaling")
	}
}
