package jwt

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
)

func Generate(ctx context.Context, t models.User) (string, error) {
	jwtSign := ctx.Value(dto.Key("jwtSign")).(string)
	key := []byte(jwtSign)

	payload := jwt.MapClaims{
		"email":         t.Email,
		"name":          t.Name,
		"surname":       t.Surname,
		"date_of_birth": t.DateOfBirth,
		"_id":           t.Id.Hex(),
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return tokenStr, err
	}

	return tokenStr, nil
}
