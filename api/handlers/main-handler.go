package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nahuelojea/handballscore/api/handlers/associations_handler"
	"github.com/nahuelojea/handballscore/api/handlers/authorization_handler"
	"github.com/nahuelojea/handballscore/api/handlers/categories_handler"
	"github.com/nahuelojea/handballscore/api/handlers/coaches_handler"
	"github.com/nahuelojea/handballscore/api/handlers/league_phase_handler"
	"github.com/nahuelojea/handballscore/api/handlers/match_coaches_handler"
	"github.com/nahuelojea/handballscore/api/handlers/match_players_handler"
	"github.com/nahuelojea/handballscore/api/handlers/matches_handler"
	"github.com/nahuelojea/handballscore/api/handlers/news_handler"
	"github.com/nahuelojea/handballscore/api/handlers/places_handler"
	"github.com/nahuelojea/handballscore/api/handlers/players_handler"
	"github.com/nahuelojea/handballscore/api/handlers/referees_handler"
	"github.com/nahuelojea/handballscore/api/handlers/teams_handler"
	"github.com/nahuelojea/handballscore/api/handlers/tournaments_category_handler"
	"github.com/nahuelojea/handballscore/api/handlers/tournaments_handler"
	"github.com/nahuelojea/handballscore/api/handlers/users_handler"
	"github.com/nahuelojea/handballscore/config/jwt"
	"github.com/nahuelojea/handballscore/dto"
)

func ProcessRequest(ctx context.Context, request events.APIGatewayProxyRequest) dto.RestResponse {

	fmt.Println("API Request: " + ctx.Value(dto.Key("path")).(string) + " > " + ctx.Value(dto.Key("method")).(string))

	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	isOk, msg, claim := jwt.ValidAuthorization(ctx, request.Headers["Authorization"])
	if !isOk {
		restResponse.Status = http.StatusUnauthorized
		restResponse.Message = msg
		return restResponse
	}
	fmt.Println("API Request Body: > " + request.Body)

	fmt.Println("User Email: > " + claim.Email)
	fmt.Println("User AssociationId: > " + claim.AssociationId)
	fmt.Println("User Role: > " + claim.Role)
	fmt.Println("User RoleId: > " + claim.RoleId)
	fmt.Println("User TeamId: > " + claim.TeamId)

	parts := strings.Split(ctx.Value(dto.Key("path")).(string), "/")

	if len(parts) > 0 {
		entityPath := parts[0]
		switch entityPath {
		case "association":
			return associations_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "auth":
			return authorization_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "category":
			return categories_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "coach":
			return coaches_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "leaguePhase":
			return league_phase_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "match":
			return matches_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "matchCoach":
			return match_coaches_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "matchPlayer":
			return match_players_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "news":
			return news_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "places":
			return places_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "player":
			return players_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "referee":
			return referees_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "team":
			return teams_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "tournamentCategory":
			return tournaments_category_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "tournament":
			return tournaments_handler.ProcessRequest(ctx, request, claim, restResponse)
		case "user":
			return users_handler.ProcessRequest(ctx, request, claim, restResponse)
		}
	}

	restResponse.Message = "Method Invalid"
	return restResponse
}
