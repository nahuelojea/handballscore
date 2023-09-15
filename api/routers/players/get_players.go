package players

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/repositories/players_repository"
	"github.com/nahuelojea/handballscore/services/categories_service"
)

func GetPlayers(request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse
	var err error

	pageStr := request.QueryStringParameters["page"]
	pageSizeStr := request.QueryStringParameters["pageSize"]
	name := request.QueryStringParameters["name"]
	surname := request.QueryStringParameters["surname"]
	dni := request.QueryStringParameters["dni"]
	gender := request.QueryStringParameters["gender"]
	teamId := request.QueryStringParameters["teamId"]
	categoryId := request.QueryStringParameters["categoryId"]
	excludeExpiredInsuranceStr := request.QueryStringParameters["excludeExpiredInsurance"]
	associationId := claim.AssociationId

	var excludeExpiredInsurance bool
	if excludeExpiredInsuranceStr != "" {
		excludeExpiredInsurance, err = strconv.ParseBool(excludeExpiredInsuranceStr)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = "'excludeExpiredInsurance' param is invalid"
			return response
		}
	}

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	var yearLimitFrom, yearLimitTo int
	if len(categoryId) > 1 {
		yearLimitFrom, yearLimitTo, gender, err = categories_service.GetLimitYearsByCategory(categoryId)
		if err != nil {
			response.Status = http.StatusNotFound
			response.Message = "Error to get category: " + err.Error()
			return response
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}

	filterOptions := players_repository.GetPlayersOptions{
		Name:                    name,
		Surname:                 surname,
		Dni:                     dni,
		Gender:                  gender,
		TeamId:                  teamId,
		ExcludeExpiredInsurance: excludeExpiredInsurance,
		YearLimitFrom:           yearLimitFrom,
		YearLimitTo:             yearLimitTo,
		AssociationId:           associationId,
		Page:                    page,
		PageSize:                pageSize,
		SortField:               "personal_data.surname",
		SortOrder:               1,
	}

	playersList, totalRecords, err := players_repository.GetPlayersFilteredAndPaginated(filterOptions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error to get players: " + err.Error()
		return response
	}

	paginatedResponse := dto.PaginatedResponse{
		TotalRecords: totalRecords,
		TotalPages:   int(totalRecords / int64(pageSize)),
		CurrentPage:  page,
		PageSize:     pageSize,
		Items:        playersList,
	}

	jsonResponse, err := json.Marshal(paginatedResponse)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formatting players to JSON: " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
