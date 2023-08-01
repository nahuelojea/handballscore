package referees_repository

import (
	"context"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	referee_collection = "referees"
)

func CreateReferee(referee models.Referee) (string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

	referee.CreatedDate = time.Now()
	referee.ModifiedDate = time.Now()
	referee.Disabled = false

	result, err := collection.InsertOne(ctx, referee)
	if err != nil {
		return "", false, err
	}

	ObjId, _ := result.InsertedID.(primitive.ObjectID)
	return ObjId.Hex(), true, nil
}

func GetReferee(ID string) (models.Referee, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

	var referee models.Referee
	objId, _ := primitive.ObjectIDFromHex(ID)

	condicion := bson.M{
		"_id": objId,
	}

	err := collection.FindOne(ctx, condicion).Decode(&referee)
	if err != nil {
		return referee, err
	}

	return referee, nil
}

func UpdateReferee(referee models.Referee, ID string) (bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

	register := make(map[string]interface{})
	if len(referee.Name) > 0 {
		register["personal_data.name"] = referee.Name
	}
	if len(referee.Surname) > 0 {
		register["personal_data.surname"] = referee.Surname
	}
	if len(referee.Avatar) > 0 {
		register["personal_data.avatar"] = referee.Avatar
	}
	register["personal_data.date_of_birth"] = referee.DateOfBirth
	if len(referee.Dni) > 0 {
		register["personal_data.dni"] = referee.Dni
	}
	if len(referee.PhoneNumber) > 0 {
		register["personal_data.phone_number"] = referee.PhoneNumber
	}
	register["status_data.modified_date"] = time.Now()

	updateString := bson.M{
		"$set": register,
	}

	objId, _ := primitive.ObjectIDFromHex(ID)
	filtro := bson.M{"_id": bson.M{"$eq": objId}}

	_, err := collection.UpdateOne(ctx, filtro, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DisableReferee(ID string) (bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

	register := make(map[string]interface{})

	register["status_data.disabled"] = true
	register["status_data.modified_date"] = time.Now()

	updateString := bson.M{
		"$set": register,
	}

	objId, _ := primitive.ObjectIDFromHex(ID)
	filtro := bson.M{"_id": bson.M{"$eq": objId}}

	_, err := collection.UpdateOne(ctx, filtro, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

func FindRefereeByDni(dni string) (models.Referee, bool, string) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

	condition := bson.M{"dni": dni}

	var result models.Referee

	err := collection.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()
	if err != nil {
		return result, false, id
	}
	return result, true, id
}
