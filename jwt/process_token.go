package jwt

import (
	"errors"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

var Email string
var UserId string

func ProcessToken(token string, JWTSign string) (*models.Claim, bool, string, error) {
	key := []byte(JWTSign)
	var claims models.Claim

	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return &claims, false, string(""), errors.New("Invalid token format")
	}

	token = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err == nil {
		// Rutina que chequea contra la BD
		_, exist, _ := users_repository.FindUserByEmail(claims.Email)
		if exist {
			Email = claims.Email
			UserId = claims.Id.Hex()
		}
		return &claims, exist, UserId, nil
	}

	if !tkn.Valid {
		return &claims, false, string(""), errors.New("Invalid token")
	}

	return &claims, false, string(""), err
}
