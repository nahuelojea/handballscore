package repositories

import (
	"context"

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

	db := MongoClient.Database(DatabaseName)
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

func FindUserByEmail(email string) (models.User, bool, string) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
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
