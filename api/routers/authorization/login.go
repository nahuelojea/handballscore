package authorization

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/dto"
	loginDTO "github.com/nahuelojea/handballscore/dto/login"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/associations_service"
	"github.com/nahuelojea/handballscore/services/authorization_service"
)

func Login(ctx context.Context) dto.RestResponse {
	var user models.User
	var response dto.RestResponse
	response.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Invalid User and/or Password " + err.Error()
		return response
	}
	if len(user.Email) == 0 {
		response.Message = "Email is required"
		return response
	}

	userData, jwtKey, refreshJwtKey, exist, err := authorization_service.Login(ctx, user.Email, user.Password)

	if !exist {
		response.Message = "Invalid User"
		return response
	}

	if err != nil {
		response.Message = err.Error()
		return response
	}

	association, _, err := associations_service.GetAssociation(userData.AssociationId)
	if err != nil {
		response.Message = "Error to get association: " + err.Error()
		return response
	}

	news := make([]loginDTO.AssociationNew, len(association.News))
	for i, n := range association.News {
		news[i] = loginDTO.AssociationNew{
			Date:  n.Date,
			Image: n.Image,
		}
	}

	resp := loginDTO.LoginResponse{
		Token:        jwtKey,
		RefreshToken: refreshJwtKey,
		Association: loginDTO.Association{
			Id:               association.Id.Hex(),
			Name:             association.Name,
			DateOfFoundation: association.DateOfFoundation,
			Email:            association.Email,
			Instagram:        association.Instagram,
			News:             news,
			Avatar:           association.Avatar,
			PhoneNumber:      association.PhoneNumber,
		},
		Role:   string(userData.Role),
		RoleId: string(userData.RoleId),
		TeamId: string(userData.TeamId),
	}

	token, err := json.Marshal(resp)
	if err != nil {
		response.Message = "Error formatting token to JSON > " + err.Error()
		return response
	}

	tokenCookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}
	cookieString := tokenCookie.String()

	refreshTokenCookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshJwtKey,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}
	refreshCookieString := refreshTokenCookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieString + "; " + refreshCookieString,
		},
	}

	response.Status = http.StatusOK
	response.Message = string(token)
	response.CustomResp = res

	return response
}
