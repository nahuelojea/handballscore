package tournaments

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/matches_service"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTournamentCategory(ctx context.Context, claim dto.Claim) dto.RestResponse {
	var tournament models.TournamentCategory
	var restResponse dto.RestResponse
	restResponse.Status = http.StatusBadRequest

	body := ctx.Value(dto.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &tournament)
	if err != nil {
		restResponse.Message = err.Error()
		return restResponse
	}

	if len(tournament.Name) == 0 {
		restResponse.Message = "Name is required"
		return restResponse
	}
	if len(tournament.CategoryId) == 0 {
		restResponse.Message = "Category id is required"
		return restResponse
	}
	if len(tournament.Teams) == 0 {
		restResponse.Message = "Teams is required"
		return restResponse
	}
	if reflect.DeepEqual(tournament.LeaguePhase, models.LeaguePhase{}) && reflect.DeepEqual(tournament.PlayoffPhase, models.PlayoffPhase{}) {
		restResponse.Message = "Type of tournament category is required"
		return restResponse
	}

	tournament.Status = models.Started

	session, err := db.MongoClient.StartSession()
	if err != nil {
		restResponse.Message = "Error starting session: " + err.Error()
		return restResponse
	}
	defer session.EndSession(context.TODO())

	err = session.StartTransaction()
	if err != nil {
		restResponse.Message = "Error starting transaction: " + err.Error()
		return restResponse
	}

	var matches []models.Match

	if !reflect.DeepEqual(tournament.LeaguePhase, models.LeaguePhase{}) {
		matches = tournament.GenerateLeagueMatches()
	}

	if !reflect.DeepEqual(tournament.PlayoffPhase, models.PlayoffPhase{}) {
		tournament.PlayoffPhase.Id = primitive.NewObjectID()
		//TODO Para pensar
	}

	id, status, err := tournaments_service.CreateTournamentCategory(claim.AssociationId, tournament)
	if err != nil {
		session.AbortTransaction(context.TODO())
		restResponse.Message = "Error to create tournament category: " + err.Error()
		return restResponse
	}
	if !status {
		session.AbortTransaction(context.TODO())
		restResponse.Message = "Error to create tournament category"
		return restResponse
	}

	_, isOk, err := matches_service.CreateMatches(claim.AssociationId, matches)
	if !isOk {
		session.AbortTransaction(context.TODO())
		restResponse.Message = "Error to create league matches"
		return restResponse
	}
	if err != nil {
		session.AbortTransaction(context.TODO())
		restResponse.Message = "Error to create league matches: " + err.Error()
		return restResponse
	}

	err = session.CommitTransaction(context.TODO())
	if err != nil {
		restResponse.Message = "Error committing transaction: " + err.Error()
		return restResponse
	}

	restResponse.Status = http.StatusCreated
	restResponse.Message = id
	return restResponse
}
