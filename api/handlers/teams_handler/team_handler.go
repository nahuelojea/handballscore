package teams_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/teams"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "team":
			return teams.AddTeam(ctx, claim)
		case "team/avatar":
			return teams.UploadAvatar(ctx, request)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "team":
			return teams.GetTeam(request)
		case "team/filter":
			return teams.GetTeams(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "team":
			return teams.UpdateTeam(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "team":
			return teams.DeleteTeam(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
