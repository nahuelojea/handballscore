package users

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/config/jwt"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

func Login(ctx context.Context) dto.RestResponse {
	var t models.User
	var r dto.RestResponse
	r.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Invalid User and/or Password " + err.Error()
		return r
	}
	if len(t.Email) == 0 {
		r.Message = "Email is required"
		return r
	}
	userData, exist := users_repository.UserLogin(t.Email, t.Password)
	if !exist {
		r.Message = "Invalid User and/or Password "
		return r
	}

	jwtKey, err := jwt.Generate(ctx, userData)
	if err != nil {
		r.Message = "Error to generate token > " + err.Error()
		return r
	}

	refreshJwtKey, err := jwt.GenerateRefreshToken(ctx, userData)
	if err != nil {
		r.Message = "Error to generate refresh token > " + err.Error()
		return r
	}

	resp := dto.LoginResponse{
		Token:         jwtKey,
		RefreshToken:  refreshJwtKey,
		Avatar:        userData.Avatar,
		AssociationId: userData.AssociationId,
	}

	token, err2 := json.Marshal(resp)
	if err2 != nil {
		r.Message = "Error formatting token to JSON > " + err2.Error()
		return r
	}

	tokenCookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}
	cookieString := tokenCookie.String()

	refreshTokenCookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshJwtKey,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}
	refreshCookieString := refreshTokenCookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieString + "; " + refreshCookieString,
		},
	}

	r.Status = http.StatusOK
	r.Message = string(token)
	r.CustomResp = res

	return r
}
