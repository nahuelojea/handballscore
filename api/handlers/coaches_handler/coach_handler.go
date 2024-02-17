package coaches_handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/routers/coaches"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim, response dto.RestResponse) dto.RestResponse {

	switch ctx.Value(dto.Key("method")).(string) {
	case "POST":
		switch ctx.Value(dto.Key("path")).(string) {
		case "coach":
			return coaches.AddCoach(ctx, claim)
		case "coach/avatar":
			return coaches.UploadAvatar(ctx, request)
		}
	case "GET":
		switch ctx.Value(dto.Key("path")).(string) {
		case "coach":
			return coaches.GetCoach(request)
		case "coach/filter":
			return coaches.GetCoachs(request, claim)
		}
	case "PUT":
		switch ctx.Value(dto.Key("path")).(string) {
		case "coach":
			return coaches.UpdateCoach(ctx, request)
		}
	case "DELETE":
		switch ctx.Value(dto.Key("path")).(string) {
		case "coach":
			return coaches.DeleteCoach(request)
		}
	}

	response.Message = "Method Invalid"
	return response
}
