package places_handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/places"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/places_service"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim models.Claim, response models.Response) models.Response {
	placeService := places_service.NewPlaceService()

	switch request.HTTPMethod {
	case "POST":
		if request.PathParameters["id"] == "" {
			return places.AddPlace(ctx, request, *placeService, claim, response)
		}
	case "GET":
		if request.PathParameters["id"] != "" {
			return places.GetPlace(ctx, request, *placeService, claim, response)
		} else {
			return places.GetPlaces(ctx, request, *placeService, claim, response)
		}
	case "PUT":
		if request.PathParameters["id"] != "" {
			return places.UpdatePlace(ctx, request, *placeService, claim, response)
		}
	case "DELETE":
		if request.PathParameters["id"] != "" {
			return places.DeletePlace(ctx, request, *placeService, claim, response)
		}
	}

	response.StatusCode = http.StatusBadRequest
	response.Message = "Invalid request method or path for places"
	return response
}
