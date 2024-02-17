package league_phase_weeks_repository

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
	league_phase_weeks_collection = "league_phase_weeks"
)

type GetLeaguePhaseWeeksOptions struct {
	LeaguePhaseId string
	Number        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateLeaguePhaseWeeks(association_id string, leaguePhaseWeeks []models.LeaguePhaseWeek) ([]string, bool, error) {
	entities := make([]models.Entity, len(leaguePhaseWeeks))
	for i, v := range leaguePhaseWeeks {
		leaguePhaseWeek := v
		entities[i] = models.Entity(&leaguePhaseWeek)
	}
	return repositories.CreateMultiple(league_phase_weeks_collection, association_id, entities)
}

func GetLeaguePhaseWeek(ID string) (models.LeaguePhaseWeek, bool, error) {
	var leaguePhaseWeek models.LeaguePhaseWeek
	_, err := repositories.GetById(league_phase_weeks_collection, ID, &leaguePhaseWeek)
	if err != nil {
		return models.LeaguePhaseWeek{}, false, err
	}

	return leaguePhaseWeek, true, nil
}

func GetLeaguePhaseWeeks(filterOptions GetLeaguePhaseWeeksOptions) ([]models.LeaguePhaseWeek, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(league_phase_weeks_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.LeaguePhaseId != "" {
		filter["league_phase_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.LeaguePhaseId, Options: "i"}}
	}
	if filterOptions.Number != "" {
		filter["number"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Number, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "number", Value: sortOrder},
		{Key: "league_phase_id", Value: sortOrder},
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

	var leaguePhaseWeeks []models.LeaguePhaseWeek
	for cur.Next(ctx) {
		var leaguePhaseWeek models.LeaguePhaseWeek
		if err := cur.Decode(&leaguePhaseWeek); err != nil {
			return nil, 0, 0, err
		}
		leaguePhaseWeeks = append(leaguePhaseWeeks, leaguePhaseWeek)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return leaguePhaseWeeks, totalRecords, totalPages, nil
}
