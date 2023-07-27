package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/jwt"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
)

func Login(ctx context.Context) models.RespApi {
	var t models.User
	var r models.RespApi
	r.Status = 400

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Invalid User and/or Password " + err.Error()
		return r
	}
	if len(t.Email) == 0 {
		r.Message = "Email is required"
		return r
	}
	userData, exist := repositories.UserLogin(t.Email, t.Password)
	if !exist {
		r.Message = "Invalid User and/or Password "
		return r
	}

	jwtKey, err := jwt.Generate(ctx, userData)
	if err != nil {
		r.Message = "Error to generate token > " + err.Error()
		return r
	}

	resp := models.LoginResponse{
		Token: jwtKey,
	}

	token, err2 := json.Marshal(resp)
	if err2 != nil {
		r.Message = "Error formatting token to JSON > " + err2.Error()
		return r
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}
	cookieString := cookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieString,
		},
	}

	r.Status = 200
	r.Message = string(token)
	r.CustomResp = res

	return r
}
