package categories_repository

import (
	"context"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	category_collection = "categories"
)

func CreateCategory(association_id string, category models.Category) (string, bool, error) {
	return repositories.Create(category_collection, association_id, &category)
}

func GetCategory(ID string) (models.Category, bool, error) {
	var category models.Category
	_, err := repositories.GetById(category_collection, ID, &category)
	if err != nil {
		return models.Category{}, false, err
	}

	return category, true, nil
}

type GetCategoriesOptions struct {
	Name          string
	Gender        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetCategoriesFilteredAndPaginated(filterOptions GetCategoriesOptions) ([]models.Category, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(category_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
	}
	if filterOptions.Gender != "" {
		filter["gender"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Gender, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

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
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var categories []models.Category
	for cur.Next(ctx) {
		var category models.Category
		if err := cur.Decode(&category); err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return categories, totalRecords, nil
}

func UpdateCategory(category models.Category, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(category.Name) > 0 {
		updateDataMap["name"] = category.Name
	}
	if len(category.Gender) > 0 {
		updateDataMap["gender"] = category.Gender
	}
	if category.AgeLimitFrom != 0 {
		updateDataMap["age_limit_from"] = category.AgeLimitFrom
	}
	if category.AgeLimitTo != 0 {
		updateDataMap["age_limit_to"] = category.AgeLimitTo
	}

	return repositories.Update(category_collection, updateDataMap, ID)
}

func DisableCategory(ID string) (bool, error) {
	return repositories.Disable(category_collection, ID)
}
