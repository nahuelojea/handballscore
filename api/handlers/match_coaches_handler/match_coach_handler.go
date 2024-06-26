package match_coaches_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/match_coaches"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {
	switch ctx.Value(dto.Key("method")).(string) {
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "matchCoach":
			return match_coaches.GetMatchCoach(request)
		}
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "matchCoach":
			return match_coaches.AddMatchCoach(ctx, claim)
		case "matchCoach/filter":
			return match_coaches.GetMatchCoaches(ctx, request, claim)
		}
	case "PATCH":
		switch ctx.Value(dto.Key("path")).(string) {
		case "matchCoach/exclusion":
			return match_coaches.UpdateExclusions(ctx, request)
		case "matchCoach/yellowCard":
			return match_coaches.UpdateYellowCard(ctx, request)
		case "matchCoach/redCard":
			return match_coaches.UpdateRedCard(ctx, request)
		case "matchCoach/blueCard":
			return match_coaches.UpdateBlueCard(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "matchCoach":
			return match_coaches.DeleteMatchCoach(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
