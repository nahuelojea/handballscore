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
			return players.AddPlayer(ctx, claim)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.GetPlayer(request)
		case "player/filter":
			return players.GetPlayers(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.UpdatePlayer(ctx, request)
		case "player/avatar":
			return players.UploadAvatar(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "player":
			return players.DeletePlayer(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
