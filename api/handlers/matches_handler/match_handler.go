package categories_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/categories"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "category":
			//return categories.AddCategory(ctx)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "category":
			return categories.GetCategory(request)
		case "category/filter":
			return categories.GetCategories(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "category":
			//return categories.UpdateCategory(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "category":
			//return categories.DisableCategory(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
