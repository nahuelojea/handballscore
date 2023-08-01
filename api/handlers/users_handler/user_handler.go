package users_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/users"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

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

	response.Message = "Method Invalid"
	return response
}
