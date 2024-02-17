package authorization_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/authorization"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "auth/login":
			return authorization.Login(ctx)
		case "auth/refresh":
			return authorization.RefreshToken(ctx, request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
