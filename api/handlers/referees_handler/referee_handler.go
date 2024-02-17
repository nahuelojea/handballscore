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
			return referees.AddReferee(ctx, claim)
		case "referee/avatar":
			return referees.UploadAvatar(ctx, request)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.GetReferee(request)
		case "referee/filter":
			return referees.GetReferees(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.UpdateReferee(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "referee":
			return referees.DeleteReferee(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
