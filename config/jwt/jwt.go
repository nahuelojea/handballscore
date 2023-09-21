package jwt

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"golang.org/x/exp/slices"
)

func getAuthorizedUrls() []string {
	return []string{
		"user/register",
		"auth/login",
		"auth/refresh"}
}

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

func ValidAuthorization(ctx context.Context, token string) (bool, string, dto.Claim) {
	path := ctx.Value(dto.Key("path")).(string)

	if slices.Contains(getAuthorizedUrls(), path) {
		return true, "", dto.Claim{}
	}

	if len(token) == 0 {
		return false, "Token is required", dto.Claim{}
	}

	claim, isOk, message, err := ProcessToken(token, ctx.Value(dto.Key("jwtSign")).(string))
	if !isOk {
		if err != nil {
			return false, err.Error(), dto.Claim{}
		} else {
			return false, message, dto.Claim{}
		}
	}

	return true, message, *claim
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
