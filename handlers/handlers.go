package handlers

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/jwt"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/routers/users"
)

func getAuthorizedUrls() []string {
	return []string{
		"user/register",
		"user/login"}
}

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest) models.RespApi {

	fmt.Println("API Request: " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var r models.RespApi
	r.Status = 400

	isOk, statusCode, msg, claim := validAuthorization(ctx, request)
	if !isOk {
		r.Status = statusCode
		r.Message = msg
		return r
	}

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "user/register":
			return users.Register(ctx)
		case "user/login":
			return users.Login(ctx)
		}
		//
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		case "user":
			return users.GetUser(request)
		}
		//
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		case "user":
			return users.UpdateUser(ctx, claim)
		}
		//
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {

		}
		//
	}

	r.Message = "Method Invalid"
	return r
}

func validAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	path := ctx.Value(models.Key("path")).(string)

	if slices.Contains(getAuthorizedUrls(), path) {
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, 401, "Token is required", models.Claim{}
	}

	claim, isOk, msg, err := jwt.ProcessToken(token, ctx.Value(models.Key("jwtSign")).(string))
	if !isOk {
		if err != nil {
			fmt.Println("Error with token " + err.Error())
			return false, 401, err.Error(), models.Claim{}
		} else {
			fmt.Println("Error with token " + msg)
			return false, 401, msg, models.Claim{}
		}
	}

	return true, 200, msg, *claim
}
