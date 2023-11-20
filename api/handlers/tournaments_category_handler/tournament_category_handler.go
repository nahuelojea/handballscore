package tournaments_category_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	tournaments "github.com/nahuelojea/handballscore/api/routers/tournaments_category"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournamentCategory":
			return tournaments.CreateTournamentCategory(ctx, claim)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournamentCategory":
			return tournaments.GetTournamentCategory(request)
		case "tournamentCategory/filter":
			return tournaments.GetTournamentsCategory(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournamentCategory":
			return tournaments.UpdateTournamentCategory(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "tournamentCategory":
			return tournaments.DeleteTournamentCategory(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
