package tournaments_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/tournaments"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournament":
			return tournaments.CreateTournament(ctx, claim)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournament":
			return tournaments.GetTournament(request)
		case "tournament/filter":
			return tournaments.GetTournaments(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournament":
			return tournaments.UpdateTournament(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournament":
			//return tournaments.DeleteTournament(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
