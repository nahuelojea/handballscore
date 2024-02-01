package users_service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
