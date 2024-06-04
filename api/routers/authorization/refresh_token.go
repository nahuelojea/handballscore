package authorization

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	loginDTO "github.com/nahuelojea/handballscore/dto/login"
	"github.com/nahuelojea/handballscore/services/authorization_service"
)

func RefreshToken(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	_, jwtKey, refreshJwtKey, err := authorization_service.RefreshToken(ctx, request.Headers["Authorization"])

	if err != nil {
		response.Message = err.Error()
		return response
	}

	resp := loginDTO.RefreshTokenResponse{
		Token:        jwtKey,
		RefreshToken: refreshJwtKey,
	}

	token, err := json.Marshal(resp)
	if err != nil {
		response.Message = "Error formatting token to JSON > " + err.Error()
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
