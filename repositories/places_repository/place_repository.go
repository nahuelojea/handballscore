package places_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	placeCollection = "places"
)

type PlaceRepository struct {
	repositories.BaseRepository[models.Place]
}

func NewPlaceRepository() *PlaceRepository {
	return &PlaceRepository{
		BaseRepository: repositories.BaseRepository[models.Place]{
			Collection: db.Database.Collection(placeCollection),
		},
	}
}

func (r *PlaceRepository) CreatePlace(ctx context.Context, place *models.Place) (string, bool, error) {
	place.SetCreatedDate()
	place.SetModifiedDate()

	result, err := r.Create(ctx, *place)
	if err != nil {
		return "", false, err
	}

	objID, ok := result.(primitive.ObjectID)
	if !ok {
		return "", false, fmt.Errorf("error converting insert id to primitive.ObjectID")
	}
	return objID.Hex(), true, nil
}

func (r *PlaceRepository) GetPlace(ctx context.Context, id string) (models.Place, bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Place{}, false, fmt.Errorf("invalid id format: %s", id)
	}
	return r.Get(ctx, objID)
}

func (r *PlaceRepository) GetPlaces(ctx context.Context, filter bson.M, page, pageSize int) ([]models.Place, int64, error) {
	findOptions := options.Find()
	if page > 0 && pageSize > 0 {
		findOptions.SetSkip(int64((page - 1) * pageSize))
		findOptions.SetLimit(int64(pageSize))
	}
	findOptions.SetSort(bson.D{{Key: "status_data.created_date", Value: -1}})

	return r.GetAll(ctx, filter, findOptions)
}

func (r *PlaceRepository) UpdatePlace(ctx context.Context, id string, place models.Place) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid id format: %s", id)
	}

	place.SetModifiedDate()
	update := bson.M{
		"$set": place,
	}

	return r.Update(ctx, objID, update)
}

func (r *PlaceRepository) DeletePlace(ctx context.Context, id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid id format: %s", id)
	}
	return r.Delete(ctx, objID)
}
