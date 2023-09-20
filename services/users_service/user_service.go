package users_service

import (
	"context"
	"errors"

	"github.com/nahuelojea/handballscore/config/jwt"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
	"github.com/nahuelojea/handballscore/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	password_encrypt_cost = 8
	AvatarUrl             = "avatars/users/"
)

func CreateUser(user models.User) (string, bool, error) {
	user.Password, _ = encryptPassword(user.Password)
	return users_repository.CreateUser(user)
}

func UserLogin(ctx context.Context, email, password string) (models.User, string, string, bool, error) {
	user, exist, _ := FindUserByEmail(email)
	if !exist {
		return user, "", "", exist, nil
	}

	passwordBytes := []byte(password)
	passwordBD := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)
	if err != nil {
		return user, "", "", exist, errors.New("Incorrect password: " + err.Error())
	}

	jwtKey, refreshJwtKey, err := jwt.GenerateTokens(ctx, user)
	if err != nil {
		return user, "", "", exist, errors.New("Error to generate tokens: " + err.Error())
	}

	return user, jwtKey, refreshJwtKey, exist, nil
}

func GetUser(ID string) (models.User, bool, error) {
	user, exist, err := users_repository.GetUser(ID)
	user.Password = ""
	return user, exist, err
}

func UpdateUser(user models.User, ID string) (bool, error) {
	return users_repository.UpdateUser(user, ID)
}

func FindUserByEmail(email string) (models.User, bool, string) {
	return users_repository.FindUserByEmail(email)
}

func encryptPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), password_encrypt_cost)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}

func UploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var user models.User

	filename = AvatarUrl + id + ".jpg"

	err := storage.UploadImage(ctx, contentType, body, filename)
	if err != nil {
		return errors.New("Error to upload image: " + err.Error())
	}

	user.SetAvatarURL(filename)
	status, err := users_repository.UpdateUser(user, id)
	if err != nil || !status {
		return errors.New("Error to update user " + err.Error())
	}
	return nil
}
