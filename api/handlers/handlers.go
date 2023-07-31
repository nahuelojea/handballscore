package handlers

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/users"
	"github.com/nahuelojea/handballscore/config/jwt"
	"github.com/nahuelojea/handballscore/dto"
)

func getAuthorizedUrls() []string {
	return []string{
		"user/register",
		"user/login"}
}

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	fmt.Println("API Request: " + ctx.Value(dto.Key("path")).(string) + " > " + ctx.Value(dto.Key("method")).(string))

	var restResponse dto.RestResponse
	restResponse.Status = 400

	isOk, statusCode, msg, claim := validAuthorization(ctx, request)
	if !isOk {
		restResponse.Status = statusCode
		restResponse.Message = msg
		return restResponse
	}

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "user/register":
			return users.Register(ctx)
		case "user/login":
			return users.Login(ctx)
		}
		//
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "user":
			return users.GetUser(request)
		}
		//
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "user":
			return users.UpdateUser(ctx, claim)
		}
		//
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {

		}
		//
	}

	restResponse.Message = "Method Invalid"
	return restResponse
}

func validAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, dto.Claim) {
	path := ctx.Value(dto.Key("path")).(string)

	if slices.Contains(getAuthorizedUrls(), path) {
		return true, 200, "", dto.Claim{}
	}

	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, 401, "Token is required", dto.Claim{}
	}

	claim, isOk, msg, err := jwt.ProcessToken(token, ctx.Value(dto.Key("jwtSign")).(string))
	if !isOk {
		if err != nil {
			fmt.Println("Error with token " + err.Error())
			return false, 401, err.Error(), dto.Claim{}
		} else {
			fmt.Println("Error with token " + msg)
			return false, 401, msg, dto.Claim{}
		}
	}

	return true, 200, msg, *claim
}
