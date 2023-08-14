package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterOptions struct {
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
	ExtraFields   map[string]interface{}
}

func Create(collectionName string, association_id string, entity models.Entity) (string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	entity.SetCreatedDate()
	entity.SetModifiedDate()
	entity.SetDisabled(false)
	entity.SetAssociationId(association_id)

	result, err := collection.InsertOne(ctx, entity)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), true, nil
}

func CreateMultiple(collectionName string, associationID string, entities []models.Entity) ([]string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	var documents []interface{}

	for _, entity := range entities {
		entity.SetCreatedDate()
		entity.SetModifiedDate()
		entity.SetDisabled(false)
		entity.SetAssociationId(associationID)
		documents = append(documents, entity)
	}

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, false, err
	}

	var insertedIDs []string
	for _, objID := range result.InsertedIDs {
		id, ok := objID.(primitive.ObjectID)
		if !ok {
			return nil, false, fmt.Errorf("Failed to get inserted ID")
		}
		insertedIDs = append(insertedIDs, id.Hex())
	}

	return insertedIDs, true, nil
}

func GetById(collectionName string, Id string, model interface{}) (interface{}, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	objId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objId}

	err = collection.FindOne(ctx, filter).Decode(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func Update(collectionName string, updateDataMap map[string]interface{}, id string) (bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	objId, _ := primitive.ObjectIDFromHex(id)

	updateDataMap["status_data.modified_date"] = time.Now()

	updateString := bson.M{
		"$set": updateDataMap,
	}

	filter := bson.M{"_id": bson.M{"$eq": objId}}

	UpdateResult, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	} else if UpdateResult.MatchedCount == 0 {
		return false, fmt.Errorf("There is no record with this id")
	}

	return true, nil
}

func Disable(collectionName string, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["status_data.disabled"] = true

	return Update(collectionName, updateDataMap, id)
}

/*func GetEntitiesFilteredAndPaginated(collectionName string, filterOptions FilterOptions) ([]models.Entity, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	for key, value := range filterOptions.ExtraFields {
		if valueStr, ok := value.(string); ok {
			filter[key] = bson.M{"$regex": primitive.Regex{Pattern: valueStr, Options: "i"}}
		}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField

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

	var entities []models.Entity
	for cur.Next(ctx) {
		var entity models.Entity
		if err := cur.Decode(&entity); err != nil {
			return nil, 0, err
		}
		entities = append(entities, entity)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return entities, totalRecords, nil
}*/
