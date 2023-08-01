package users_repository

import (
	"context"
	"fmt"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	user_collection       = "users"
	password_encrypt_cost = 8
)

func CreateUser(u models.User) (string, bool, error) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(user_collection)

	u.Password, _ = encryptPassword(u.Password)

	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return "", false, err
	}

	ObjId, _ := result.InsertedID.(primitive.ObjectID)
	return ObjId.String(), true, nil
}

func UserLogin(email string, password string) (models.User, bool) {
	usu, encontrado, _ := FindUserByEmail(email)
	if !encontrado {
		return usu, false
	}

	passwordBytes := []byte(password)
	passwordBD := []byte(usu.Password)

	err := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)
	if err != nil {
		return usu, false
	}

	return usu, true
}

func GetUser(ID string) (models.User, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(user_collection)

	var user models.User
	objId, _ := primitive.ObjectIDFromHex(ID)

	condicion := bson.M{
		"_id": objId,
	}

	err := collection.FindOne(ctx, condicion).Decode(&user)
	user.Password = ""
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(user models.User, ID string) (bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(user_collection)

	register := make(map[string]interface{})
	if len(user.Name) > 0 {
		register["personal_data.name"] = user.Name
	}
	if len(user.Surname) > 0 {
		register["personal_data.surname"] = user.Surname
	}
	if len(user.Avatar) > 0 {
		register["personal_data.avatar"] = user.Avatar
	}
	register["personal_data.date_of_birth"] = user.DateOfBirth
	if len(user.Dni) > 0 {
		register["personal_data.dni"] = user.Dni
	}
	if len(user.PhoneNumber) > 0 {
		register["personal_data.phone_number"] = user.PhoneNumber
	}

	updateString := bson.M{
		"$set": register,
	}

	objId, _ := primitive.ObjectIDFromHex(ID)
	filtro := bson.M{"_id": bson.M{"$eq": objId}}

	fmt.Println("ID a modificar: ", ID, objId)

	_, err := collection.UpdateOne(ctx, filtro, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

func FindUserByEmail(email string) (models.User, bool, string) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(user_collection)

	condition := bson.M{"email": email}

	var result models.User

	err := collection.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()
	if err != nil {
		return result, false, id
	}
	return result, true, id
}

func encryptPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), password_encrypt_cost)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}
