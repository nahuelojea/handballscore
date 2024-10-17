package news_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	/*switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "new":
			return news.AddNew(ctx, claim)
		case "new/image":
			return news.UploadImage(ctx, request)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "new":
			return news.GetNew(request)
		case "new/filter":
			return news.GetNews(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "new":
			return news.UpdateNew(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "new":
			return news.DeleteNew(request)
		}
	}*/

	response.Message = "Method Invalid"
	return response
}
