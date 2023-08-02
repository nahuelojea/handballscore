package referees

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
)

func GetReferees(request events.APIGatewayProxyRequest) dto.RestResponse {
	var response dto.RestResponse

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	surname := request.QueryStringParameters["surname"]
	dni := request.QueryStringParameters["dni"]
	associationId := request.QueryStringParameters["associationId"]

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	// Convertir los parámetros de paginación a números enteros
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // Si no se proporciona o es inválido, usar página 1 por defecto
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20 // Si no se proporciona o es inválido, usar tamaño de página 10 por defecto
	}

	filterOptions := referees_repository.GetRefereesOptions{
		Name:          name,
		Surname:       surname,
		Dni:           dni,
		AssociationId: associationId,
		Page:          page,
		PageSize:      pageSize,
		SortField:     "personal_data.surname",
		SortOrder:     1,
	}

	refereesList, totalRecords, err := referees_repository.GetRefereesFilteredAndPaginated(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get referees: " + err.Error()
		return response
	}

	// Crear una estructura para la respuesta paginada
	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   int(totalRecords / int64(pageSize)),
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        refereesList,
	}

	// Convertir la respuesta paginada a JSON
	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting referees to JSON: " + err.Error()
		return response
	}

	// Asignar la respuesta JSON a la estructura de la respuesta del APIGateway
	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
