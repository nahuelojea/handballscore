package associations

import (
	"encoding/json"
	"net/http"

	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/services/associations_service"
)

func GetAssociation(claim dto.Claim) dto.RestResponse {
	var response dto.RestResponse

	associationId := claim.AssociationId

	if len(associationId) < 1 {
		response.Status = http.StatusBadRequest
		response.Message = "'associationId' param is mandatory"
		return response
	}

	association, _, err := associations_service.GetAssociation(associationId)
	if err != nil {
		response.Status = http.StatusNotFound
		response.Message = "Error to get association: " + err.Error()
		return response
	}

	jsonResponse, err := json.Marshal(association)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error formating association to JSON " + err.Error()
		return response
	}

	response.Status = http.StatusOK
	response.Message = string(jsonResponse)
	return response
}
