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

func Create(collectionName string, entity models.Entity) (string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	entity.SetCreatedDate()
	entity.SetModifiedDate()
	entity.SetDisabled(false)

	result, err := collection.InsertOne(ctx, entity)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), true, nil
}

func GetById(collectionName string, Id string) (models.Entity, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(collectionName)

	var entity models.Entity
	objId, _ := primitive.ObjectIDFromHex(Id)

	filter := bson.M{
		"_id": objId,
	}

	err := collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
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
