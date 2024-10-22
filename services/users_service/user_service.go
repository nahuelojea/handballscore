package users_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
	"github.com/nahuelojea/handballscore/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordEncryptCost = 8
	AvatarUrl             = "avatars/users/"
	standardPassword	  = "123"
)

func CreateUser(user models.User) (string, bool, error) {
	user.Password, _ = encryptPassword(user.Password)
	return users_repository.CreateUser(user)
}

func GetUser(ID string) (models.User, bool, error) {
	user, exist, err := users_repository.GetUser(ID)
	user.Password = ""
	return user, exist, err
}

func GetUserWithPassword(ID string) (models.User, bool, error) {
	user, exist, err := users_repository.GetUser(ID)
	return user, exist, err
}

func UpdateUser(user models.User, ID string) (bool, error) {
	return users_repository.UpdateUser(user, ID)
}

func FindUserByEmail(email string) (models.User, bool, string) {
	return users_repository.FindUserByEmail(email)
}

func encryptPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), passwordEncryptCost)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}

func UploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var user models.User

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename = fmt.Sprintf("%s%d_%s.jpg", AvatarUrl, timestamp, id)

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

func ChangePassword(id, oldPassword, newPassword string) error {
	user, exist, _ := GetUserWithPassword(id)
	if !exist {
		return errors.New("User does not exist")
	}

	oldPasswordBytes := []byte(oldPassword)
	passwordBD := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(passwordBD, oldPasswordBytes)
	if err != nil {
		return errors.New("Incorrect old password: " + err.Error())
	}

	encryptedPassword, err := encryptPassword(newPassword)
	if err != nil {
		return errors.New("Error encrypting password: " + err.Error())
	}

	user.Password = encryptedPassword
	status, err := users_repository.UpdateUser(user, user.Id.Hex())
	if err != nil || !status {
		return errors.New("Error updating user password: " + err.Error())
	}
	return nil
}

func ResetPassword(id string) error {
	user, exist, _ := GetUser(id)
	if !exist {
		return errors.New("User does not exist")
	}

	emailPrefix := user.Email[:strings.Index(user.Email, "@")]
	newPassword := emailPrefix + standardPassword

	encryptedPassword, err := encryptPassword(newPassword)
	if err != nil {
		return errors.New("Error encrypting password: " + err.Error())
	}

	user.Password = encryptedPassword
	status, err := users_repository.UpdateUser(user, user.Id.Hex())
	if err != nil || !status {
		return errors.New("Error updating user password: " + err.Error())
	}
	return nil
}