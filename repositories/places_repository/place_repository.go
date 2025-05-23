package places_repository

import (
	"context"
	"math"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const place_collection = "places"

type GetPlacesOptions struct {
	Name          string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreatePlace(association_id string, place models.Place) (string, bool, error) {
	return repositories.Create(place_collection, association_id, &place)
}

func GetPlace(ID string) (models.Place, bool, error) {
	var place models.Place
	_, err := repositories.GetById(place_collection, ID, &place)
	if err != nil {
		return models.Place{}, false, err
	}
	return place, true, nil
}

func GetPlaces(filterOptions GetPlacesOptions) ([]models.Place, int64, int, error) {
	ctx := context.TODO()
	mongoDB := db.MongoClient.Database(db.DatabaseName)
	collection := mongoDB.Collection(place_collection)

	filter := bson.M{}

	if len(filterOptions.AssociationId) > 0 {
		filter["association_id"] = filterOptions.AssociationId
	}
	if len(filterOptions.Name) > 0 {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
	}

	page := filterOptions.Page
	if page == 0 {
		page = 1
	}
	pageSize := filterOptions.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "name"
	}
	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var places []models.Place
	for cur.Next(ctx) {
		var place models.Place
		if err := cur.Decode(&place); err != nil {
			return nil, 0, 0, err
		}
		places = append(places, place)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	}


	return places, totalRecords, totalPages, nil
}

func UpdatePlace(place models.Place, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	if len(place.Name) > 0 {
		updateDataMap["name"] = place.Name
	}

	if place.Ubication.Latitude != 0 {
		updateDataMap["ubication.latitude"] = place.Ubication.Latitude
	}
	if place.Ubication.Longitude != 0 {
		updateDataMap["ubication.longitude"] = place.Ubication.Longitude
	}
	
	if len(place.AssociationId) > 0 {
	    updateDataMap["association_id"] = place.AssociationId
	}


	return repositories.Update(place_collection, updateDataMap, ID)
}

func DeletePlace(ID string) (bool, error) {
	return repositories.Delete(place_collection, ID)
}
