package match_players_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/match_players"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {
	switch ctx.Value(dto.Key("method")).(string) {
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "matchPlayer":
			return match_players.GetMatchPlayer(request)
		case "matchPlayer/filter":
			return match_players.GetMatchPlayers(ctx, request, claim)
		}
	case "PATCH":
		switch ctx.Value(dto.Key("path")).(string) {
		case "matchPlayer/goal":
			return match_players.UpdateGoal(ctx, request)
		case "matchPlayer/exclusion":
			return match_players.UpdateExclusions(ctx, request)
		case "matchPlayer/yellowCard":
			return match_players.UpdateYellowCard(ctx, request)
		case "matchPlayer/redCard":
			return match_players.UpdateRedCard(ctx, request)
		case "matchPlayer/blueCard":
			return match_players.UpdateBlueCard(ctx, request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
