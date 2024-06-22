package players

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/players_service"
)

func ImportFromExcel(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	fileContent := request.Body

	if len(fileContent) == 0 {
		response.Message = "File content is empty"
		return response
	}

	_, errors := players_service.ImportFromExcel(fileContent)
	if errors != nil {
		response.Message = "Error to import file"
		return response
	}

	response.Message = "File processed successfully"
	return response
}
