package jwt

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
)

func GenerateTokens(ctx context.Context, user models.User) (string, string, error) {
	jwtSign := ctx.Value(dto.Key("jwtSign")).(string)
	key := []byte(jwtSign)

	token := generateToken(user)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}

	refreshToken := generateRefreshToken(user)
	refreshTokenStr, err := refreshToken.SignedString(key)
	if err != nil {
		return "", "", err
	}

	return tokenStr, refreshTokenStr, nil
}

func generateToken(user models.User) *jwt.Token {
	payload := jwt.MapClaims{
		"email":          user.Email,
		"role":           user.Role,
		"association_id": user.AssociationId,
		"_id":            user.Id.Hex(),
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token
}

func generateRefreshToken(user models.User) *jwt.Token {
	payload := jwt.MapClaims{
		"email":          user.Email,
		"role":           user.Role,
		"association_id": user.AssociationId,
		"_id":            user.Id.Hex(),
		"exp":            time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token
}
