package league_phases_repository

import (
	"context"
	"math"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	league_phase_collection = "league_phases"
)

type GetLeaguePhasesOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func CreateLeaguePhase(association_id string, leaguePhase models.LeaguePhase) (string, bool, error) {
	return repositories.Create(league_phase_collection, association_id, &leaguePhase)
}

func GetLeaguePhase(ID string) (models.LeaguePhase, bool, error) {
	var leaguePhase models.LeaguePhase
	_, err := repositories.GetById(league_phase_collection, ID, &leaguePhase)
	if err != nil {
		return models.LeaguePhase{}, false, err
	}

	return leaguePhase, true, nil
}

func GetLeaguePhases(filterOptions GetLeaguePhasesOptions) ([]models.LeaguePhase, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(league_phase_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.TournamentCategoryId != "" {
		filter["tournament_category_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.TournamentCategoryId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "tournament_category_id", Value: sortOrder},
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(sortFields)

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var leaguePhases []models.LeaguePhase
	for cur.Next(ctx) {
		var leaguePhase models.LeaguePhase
		if err := cur.Decode(&leaguePhase); err != nil {
			return nil, 0, 0, err
		}
		leaguePhases = append(leaguePhases, leaguePhase)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return leaguePhases, totalRecords, totalPages, nil
}

func DeleteLeaguePhase(ID string) (bool, error) {
	return repositories.Delete(league_phase_collection, ID)
}
