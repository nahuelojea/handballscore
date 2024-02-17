package associations_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/associations"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "association":
			return associations.GetAssociation(claim)
		case "association/filter":
			return associations.GetAssociations(request, claim)
		}
	}

	response.Message = "Method Invalid"
	return response
}
