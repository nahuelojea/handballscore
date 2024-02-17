package authorization_service

import (
	"context"
	"errors"

	"github.com/nahuelojea/handballscore/config/jwt"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/users_service"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, email, password string) (models.User, string, string, bool, error) {
	user, exist, _ := users_service.FindUserByEmail(email)
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

func RefreshToken(ctx context.Context, refreshToken string) (models.User, string, string, error) {
	var user models.User

	if len(refreshToken) == 0 {
		return user, "", "", errors.New("Refresh token is required")
	}

	claim, isOk, message, err := jwt.ProcessToken(refreshToken, ctx.Value(dto.Key("jwtSign")).(string))
	if !isOk {
		if err != nil {
			return user, "", "", err
		} else {
			return user, "", "", errors.New("Error with token: " + message)
		}
	}

	user, exist, _ := users_service.FindUserByEmail(claim.Email)
	if !exist {
		return user, "", "", nil
	}

	jwtKey, refreshJwtKey, err := jwt.GenerateTokens(ctx, user)
	if err != nil {
		return user, "", "", errors.New("Error to generate tokens: " + err.Error())
	}

	return user, jwtKey, refreshJwtKey, nil
}
