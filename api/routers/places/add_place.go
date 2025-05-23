package places

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/places_service"
)

func AddPlace(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	response.Message = "Error to add place"
	response.Status = http.StatusBadRequest

	var place models.Place
	body := request.Body
	err := json.Unmarshal([]byte(body), &place)
	if err != nil {
		response.Message = "Error to parse place: " + err.Error()
		return response
	}

	if len(place.Name) == 0 {
		response.Message = "Name is required"
		return response
	}

	place.SetAssociationId(claim.AssociationId)

	id, status, err := places_service.CreatePlace(ctx, &place)
	if err != nil {
		response.Message = "Error to create place: " + err.Error()
		return response
	}

	if !status {
		response.Message = "Error to insert place in database"
		return response
	}

	response.StatusCode = http.StatusCreated
	response.Message = "Place created"
	response.Data = dto.CreateIdResponse{Id: id}
	return response
}
