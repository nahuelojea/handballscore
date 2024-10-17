package league_phase_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/league_phase"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {
	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "leaguePhase/recalculateTeamScores":
			return league_phase.RecalculateLeaguePhase(request)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {

		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {

		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {

		}
	}

	response.Message = "Method Invalid"
	return response
}
