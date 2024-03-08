package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	categories_dto "github.com/nahuelojea/handballscore/dto/tournament_categories"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func GetWeeksAndRounds(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var weeksAndRoundsResponse []categories_dto.WeeksAndRoundsResponse

	id := request.QueryStringParameters["id"]
	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	associationId := claim.AssociationId

	if len(associationId) < 1 {
		return dto.RestResponse{
			Status:  http.StatusBadRequest,
			Message: "'associationId' is mandatory",
		}
	}
	if len(id) < 1 {
		return dto.RestResponse{
			Status:  http.StatusBadRequest,
			Message: "'id' param is mandatory",
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	weeksAndRoundsResponse, totalRecords, totalPages, err := tournaments_service.GetWeeksAndRounds(id, associationId, page, pageSize)
	if err != nil {
		return dto.RestResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error to get weeks and rounds: " + err.Error(),
		}
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        weeksAndRoundsResponse,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		return dto.RestResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error formatting categories to JSON: " + err.Error(),
		}
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
