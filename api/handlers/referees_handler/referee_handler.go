package referees_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/referees"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.AddReferee(ctx)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.GetReferee(request)
		case "referee/filter":
			return referees.GetReferees(request)
		case "referee/avatar":
			return referees.GetAvatar(ctx, request)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.UpdateReferee(ctx, request)
		case "referee/avatar":
			return referees.UpdateAvatar(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.DisableReferee(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
