package users_repository

import (
	"context"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	user_collection = "users"
)

func CreateUser(user models.User) (string, bool, error) {
	return repositories.Create(user_collection, user.AssociationId, &user)
}

func GetUser(ID string) (models.User, bool, error) {
	var user models.User
	_, err := repositories.GetById(user_collection, ID, &user)
	if err != nil {
		return models.User{}, false, err
	}
	return user, true, nil
}

func UpdateUser(user models.User, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(user.Avatar) > 0 {
		updateDataMap["avatar"] = user.Avatar
	}

	return repositories.Update(user_collection, updateDataMap, ID)
}

func FindUserByEmail(email string) (models.User, bool, string) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(user_collection)

	filter := bson.M{"email": email}

	var result models.User

	err := collection.FindOne(ctx, filter).Decode(&result)
	id := result.Id.Hex()
	if err != nil {
		return result, false, id
	}
	return result, true, id
}
