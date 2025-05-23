package places

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/places_service"
)

func GetPlace(ctx context.Context, request events.APIGatewayProxyRequest, service places_service.PlaceService, claim models.Claim, response models.Response) models.Response {
	response.Message = "Error to get place"
	response.StatusCode = http.StatusBadRequest

	placeId := request.PathParameters["id"]
	if len(placeId) == 0 {
		response.Message = "Id is required"
		return response
	}

	place, status, err := service.GetPlace(ctx, placeId)
	if err != nil {
		response.Message = "Error to get place: " + err.Error()
		return response
	}

	if !status {
		response.StatusCode = http.StatusNotFound
		response.Message = "Place not found"
		return response
	}

	if place.AssociationId != claim.AssociationId {
		response.StatusCode = http.StatusForbidden
		response.Message = "User does not have permission to access this place"
		return response
	}

	response.StatusCode = http.StatusOK
	response.Message = "Place"
	response.Data = place
	return response
}
