package places

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/places_service"
)

func UpdatePlace(ctx context.Context, request events.APIGatewayProxyRequest, service places_service.PlaceService, claim models.Claim, response models.Response) models.Response {
	response.Message = "Error to update place"
	response.StatusCode = http.StatusBadRequest

	placeId := request.PathParameters["id"]
	if len(placeId) == 0 {
		response.Message = "Id is required"
		return response
	}

	var place models.Place
	body := request.Body
	err := json.Unmarshal([]byte(body), &place)
	if err != nil {
		response.Message = "Error to parse place: " + err.Error()
		return response
	}

	existingPlace, status, err := service.GetPlace(ctx, placeId)
	if err != nil {
		response.Message = "Error to get place for validation: " + err.Error()
		return response
	}
	if !status {
		response.StatusCode = http.StatusNotFound
		response.Message = "Place not found for validation"
		return response
	}
	if existingPlace.AssociationId != claim.AssociationId {
		response.StatusCode = http.StatusForbidden
		response.Message = "User does not have permission to update this place"
		return response
	}

	place.SetAssociationId(claim.AssociationId)
	place.SetId(existingPlace.Id)

	updated, err := service.UpdatePlace(ctx, placeId, place)
	if err != nil {
		response.Message = "Error to update place: " + err.Error()
		return response
	}

	if !updated {
		response.Message = "Error to update place in database"
		return response
	}

	response.StatusCode = http.StatusOK
	response.Message = "Place updated"
	return response
}
