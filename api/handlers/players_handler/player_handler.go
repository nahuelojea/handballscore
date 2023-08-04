package players_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/players"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.AddPlayer(ctx)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.GetPlayer(request)
		case "player/filter":
			return players.GetPlayers(request)
		case "player/avatar":
			return players.GetAvatar(ctx, request)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.UpdatePlayer(ctx, request)
		case "player/avatar":
			return players.UpdateAvatar(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.DisablePlayer(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
