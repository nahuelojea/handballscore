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

	if len(user.Password) > 0 {
		updateDataMap["password"] = user.Password
	}

	if len(user.Name) > 0 {
		updateDataMap["personal_data.name"] = user.Name
	}
	if len(user.Surname) > 0 {
		updateDataMap["personal_data.surname"] = user.Surname
	}
	if len(user.Dni) > 0 {
		updateDataMap["personal_data.dni"] = user.Dni
	}
	if !user.DateOfBirth.IsZero() {
		updateDataMap["personal_data.date_of_birth"] = user.DateOfBirth
	}
	if len(user.Gender) > 0 {
		updateDataMap["personal_data.gender"] = user.Gender
	}
	if len(user.PhoneNumber) > 0 {
		updateDataMap["personal_data.phone_number"] = user.PhoneNumber
	}
	updateDataMap["personal_data.disabled"] = user.Disabled

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

func UpdateAvatar(user models.User, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(user.Avatar) > 0 {
		updateDataMap["personal_data.avatar"] = user.Avatar
	}

	return repositories.Update(user_collection, updateDataMap, ID)
}
