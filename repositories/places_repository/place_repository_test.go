package places_repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestPlaceRepository_CreatePlace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		placeRepo := NewPlaceRepository()
		placeRepo.Collection = mt.Coll
		expectedID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "insertedID", Value: expectedID}))

		place := &models.Place{
			Name:          "Test Place",
			AssociationId: "assoc1",
		}

		id, status, err := placeRepo.CreatePlace(context.Background(), place)

		assert.NoError(t, err)
		assert.True(t, status)
		assert.Equal(t, expectedID.Hex(), id)
		assert.NotZero(t, place.Status_Data.CreatedDate)
		assert.NotZero(t, place.Status_Data.ModifiedDate)
	})

	mt.Run("failure - insert error", func(mt *mtest.T) {
		placeRepo := NewPlaceRepository()
		placeRepo.Collection = mt.Coll

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		place := &models.Place{Name: "Test Place"}
		_, status, err := placeRepo.CreatePlace(context.Background(), place)

		assert.Error(t, err)
		assert.False(t, status)
	})
}

func TestPlaceRepository_GetPlace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	placeRepo := NewPlaceRepository()

	expectedID := primitive.NewObjectID()
	place := models.Place{
		Id:            expectedID,
		Name:          "Found Place",
		AssociationId: "assoc1",
	}

	mt.Run("success", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", mt.DB.Name(), mt.Coll.Name()), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedID},
			{Key: "name", Value: place.Name},
			{Key: "association_id", Value: place.AssociationId},
		}))

		foundPlace, status, err := placeRepo.GetPlace(context.Background(), expectedID.Hex())

		assert.NoError(t, err)
		assert.True(t, status)
		assert.Equal(t, expectedID, foundPlace.Id)
		assert.Equal(t, place.Name, foundPlace.Name)
	})

	mt.Run("not found", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateCursorResponse(0, fmt.Sprintf("%s.%s", mt.DB.Name(), mt.Coll.Name()), mtest.FirstBatch))

		_, status, err := placeRepo.GetPlace(context.Background(), primitive.NewObjectID().Hex())

		assert.Error(t, err) // Should be an error like "mongo: no documents in result"
		assert.False(t, status)
	})

	mt.Run("invalid id format", func(mt *mtest.T) {
		_, status, err := placeRepo.GetPlace(context.Background(), "invalidID")
		assert.Error(t, err)
		assert.False(t, status)
		assert.Contains(t, err.Error(), "invalid id format")
	})
}

func TestPlaceRepository_GetPlaces(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	placeRepo := NewPlaceRepository()

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	place1 := bson.D{
		{Key: "_id", Value: id1}, {Key: "name", Value: "Place Alpha"}, {Key: "association_id", Value: "assoc1"},
		{Key: "status_data", Value: bson.D{{Key: "created_date", Value: time.Now().Add(-time.Hour)}}},
	}
	place2 := bson.D{
		{Key: "_id", Value: id2}, {Key: "name", Value: "Place Beta"}, {Key: "association_id", Value: "assoc1"},
		{Key: "status_data", Value: bson.D{{Key: "created_date", Value: time.Now()}}},
	}

	mt.Run("success - no filter, with pagination", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll

		findResponse := mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", mt.DB.Name(), mt.Coll.Name()), mtest.FirstBatch, place2)
		countResponse := mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 2}}...)
		mt.AddMockResponses(findResponse, countResponse)

		places, total, err := placeRepo.GetPlaces(context.Background(), bson.M{"association_id": "assoc1"}, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, places, 1)
		assert.Equal(t, id2, places[0].Id) // Sorted by created_date desc
	})

	mt.Run("success - with name filter", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		findResponse := mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", mt.DB.Name(), mt.Coll.Name()), mtest.FirstBatch, place1)
		countResponse := mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...)
		mt.AddMockResponses(findResponse, countResponse)

		filter := bson.M{"name": "Alpha", "association_id": "assoc1"}
		places, total, err := placeRepo.GetPlaces(context.Background(), filter, 1, 5)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, places, 1)
		assert.Equal(t, "Place Alpha", places[0].Name)
	})
}

func TestPlaceRepository_UpdatePlace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	placeRepo := NewPlaceRepository()
	placeID := primitive.NewObjectID()

	mt.Run("success", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "nModified", Value: 1}))

		placeToUpdate := models.Place{
			Name: "Updated Name",
			Ubication: models.UbicationCoordinates{
				Latitude:  1.0,
				Longitude: 1.0,
			},
		}
		status, err := placeRepo.UpdatePlace(context.Background(), placeID.Hex(), placeToUpdate)

		assert.NoError(t, err)
		assert.True(t, status)
		// ModifiedDate should be updated by the method
		// We can't directly check placeToUpdate.Status_Data.ModifiedDate here
		// as it's updated within the method on a copy, but the $set operation would include it.
	})

	mt.Run("not found", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "nModified", Value: 0})) // Simulate no document matched

		status, err := placeRepo.UpdatePlace(context.Background(), placeID.Hex(), models.Place{Name: "Any"})
		assert.NoError(t, err) // Update itself doesn't error if no doc matches
		assert.False(t, status) // Status indicates no document was updated
	})

	mt.Run("invalid id format", func(mt *mtest.T) {
		_, err := placeRepo.UpdatePlace(context.Background(), "invalidID", models.Place{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id format")
	})
}

func TestPlaceRepository_DeletePlace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	placeRepo := NewPlaceRepository()
	placeID := primitive.NewObjectID()

	mt.Run("success", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1})) // n indicates documents deleted

		status, err := placeRepo.DeletePlace(context.Background(), placeID.Hex())
		assert.NoError(t, err)
		assert.True(t, status)
	})

	mt.Run("not found", func(mt *mtest.T) {
		placeRepo.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 0})) // Simulate no document matched

		status, err := placeRepo.DeletePlace(context.Background(), placeID.Hex())
		assert.NoError(t, err)  // Delete itself doesn't error if no doc matches
		assert.False(t, status) // Status indicates no document was deleted
	})

	mt.Run("invalid id format", func(mt *mtest.T) {
		_, err := placeRepo.DeletePlace(context.Background(), "invalidID")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id format")
	})
}
