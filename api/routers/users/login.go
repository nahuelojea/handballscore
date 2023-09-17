package users

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/users_service"
)

func Login(ctx context.Context) dto.RestResponse {
	var user models.User
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Invalid User and/or Password " + err.Error()
		return response
	}
	if len(user.Email) == 0 {
		response.Message = "Email is required"
		return response
	}

	userData, jwtKey, refreshJwtKey, exist, err := users_service.UserLogin(ctx, user.Email, user.Password)

	if !exist {
		response.Message = "Invalid User"
		return response
	}

	if err != nil {
		response.Message = err.Error()
		return response
	}

	resp := dto.LoginResponse{
		Token:         jwtKey,
		RefreshToken:  refreshJwtKey,
		Avatar:        userData.Avatar,
		AssociationId: userData.AssociationId,
	}

	token, err2 := json.Marshal(resp)
	if err2 != nil {
		response.Message = "Error formatting token to JSON > " + err2.Error()
		return response
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

	response.Status = http.StatusOK
	response.Message = string(token)
	response.CustomResp = res

	return response
}
