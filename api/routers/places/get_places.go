package places

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/places_service"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPlaces(ctx context.Context, request events.APIGatewayProxyRequest, service places_service.PlaceService, claim models.Claim, response models.Response) models.Response {
	response.Message = "Error to get places"
	response.StatusCode = http.StatusBadRequest

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]

	page, _ := strconv.Atoi(pageStr)
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize == 0 {
		pageSize = 20
	}

	filter := bson.M{"association_id": claim.AssociationId}
	if len(name) > 0 {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	places, totalRecords, err := service.GetPlaces(ctx, filter, page, pageSize)
	if err != nil {
		response.Message = "Error to get places: " + err.Error()
		return response
	}

	response.StatusCode = http.StatusOK
	response.Message = "Places list"
	response.Total = totalRecords
	response.Data = places
	return response
}
