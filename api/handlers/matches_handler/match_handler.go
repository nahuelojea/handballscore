package matches_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/matches"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "match":
			return matches.GetMatch(request)
		case "match/header":
			return matches.GetMatchHeader(request)
		case "match/filter":
			return matches.GetMatches(request, claim)
		case "match/journey":
			return matches.GetMatchesByJourney(request, claim)
		case "match/today":
			return matches.GetMatchesToday(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "match/program":
			return matches.ProgramMatch(ctx, request)
		case "match/start":
			return matches.StartMatch(ctx, request)
		}
	case "PATCH":
		switch ctx.Value(dto.Key("path")).(string) {
		case "match/timeout":
			return matches.UpdateTimeouts(ctx, request)
		case "match/startSecondHalf":
			return matches.StartSecondHalf(ctx, request)
		case "match/end":
			return matches.EndMatch(ctx, request)
		case "match/suspend":
			return matches.SuspendMatch(ctx, request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
