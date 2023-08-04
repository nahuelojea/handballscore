package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/handlers/coaches_handler"
	"github.com/nahuelojea/handballscore/api/handlers/players_handler"
	"github.com/nahuelojea/handballscore/api/handlers/referees_handler"
	"github.com/nahuelojea/handballscore/api/handlers/teams_handler"
	"github.com/nahuelojea/handballscore/api/handlers/users_handler"
	"github.com/nahuelojea/handballscore/config/jwt"
	"github.com/nahuelojea/handballscore/dto"
)

func getAuthorizedUrls() []string {
	return []string{
		"user/register",
		"user/login"}
}

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	fmt.Println("API Request: " + ctx.Value(dto.Key("path")).(string) + " > " + ctx.Value(dto.Key("method")).(string))

	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	isOk, statusCode, msg, claim := validAuthorization(ctx, request)
	if !isOk {
		restResponse.Status = statusCode
		restResponse.Message = msg
		return restResponse
	}

	parts := strings.Split(ctx.Value(dto.Key("path")).(string), "/")

	if len(parts) > 0 {
		entityPath := parts[0]

		switch entityPath {
		case "coach":
			return coaches_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "user":
			return users_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "player":
			return players_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "team":
			return teams_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "referee":
			return referees_handler.ProcessRequest(ctx, request, claim, restResponse)
		}
	}

	restResponse.Message = "Method Invalid"
	return restResponse
}

func validAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, dto.Claim) {
	path := ctx.Value(dto.Key("path")).(string)

	if slices.Contains(getAuthorizedUrls(), path) {
		return true, http.StatusOK, "", dto.Claim{}
	}

	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, http.StatusUnauthorized, "Token is required", dto.Claim{}
	}

	claim, isOk, msg, err := jwt.ProcessToken(token, ctx.Value(dto.Key("jwtSign")).(string))
	if !isOk {
		if err != nil {
			fmt.Println("Error with token " + err.Error())
			return false, http.StatusUnauthorized, err.Error(), dto.Claim{}
		} else {
			fmt.Println("Error with token " + msg)
			return false, http.StatusUnauthorized, msg, dto.Claim{}
		}
	}

	return true, http.StatusOK, msg, *claim
}
