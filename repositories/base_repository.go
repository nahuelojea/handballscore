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

func Create(collectionName, associationId string, entity models.Entity) (string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	entity.SetCreatedDate()
	entity.SetModifiedDate()
	entity.SetAssociationId(associationId)

	result, err := collection.InsertOne(ctx, entity)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), true, nil
}

func CreateMultiple(collectionName, associationId string, entities []models.Entity) ([]string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	var documents []interface{}

	for _, entity := range entities {
		entity.SetCreatedDate()
		entity.SetModifiedDate()
		entity.SetAssociationId(associationId)
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

func GetById(collectionName, Id string, model interface{}) (interface{}, error) {
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

func FindOne(collectionName, associationID string, condition bson.M, model interface{}) (interface{}, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	baseCondition := bson.M{
		"association_id": associationID,
	}

	combinedCondition := bson.M{
		"$and": []bson.M{
			baseCondition,
			condition,
		},
	}

	err := collection.FindOne(ctx, combinedCondition).Decode(model)
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

func Delete(collectionName string, id string) (bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": objId}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return true, nil
}
