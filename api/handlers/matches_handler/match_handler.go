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
			//return categories.GetCategory(request)
		case "match/filter":
			//return categories.GetCategories(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "match/program":
			return matches.ProgramMatch(ctx, request)
		case "match/start":
			return matches.StartMatch(ctx, request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
