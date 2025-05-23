package places_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/places"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "place":
			return places.AddPlace(ctx, claim)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "place":
			return places.GetPlace(request)
		case "place/filter":
			return places.GetPlaces(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "place":
			return places.UpdatePlace(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "place":
			return places.DeletePlace(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
