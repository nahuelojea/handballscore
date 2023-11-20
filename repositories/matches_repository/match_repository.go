package matches_repository

import (
	"context"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	match_collection = "matches"
)

func CreateMatches(associationID string, matches []models.Match) ([]string, bool, error) {
	entities := make([]models.Entity, len(matches))
	for i, v := range matches {
		match := v
		entities[i] = models.Entity(&match)
	}

	return repositories.CreateMultiple(match_collection, associationID, entities)
}

func CreateMatch(association_id string, match models.Match) (string, bool, error) {
	return repositories.Create(match_collection, association_id, &match)
}

func GetMatch(ID string) (models.Match, bool, error) {
	var match models.Match
	_, err := repositories.GetById(match_collection, ID, &match)
	if err != nil {
		return models.Match{}, false, err
	}

	return match, true, nil
}

type GetMatchesOptions struct {
	PhaseId       string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetMatches(filterOptions GetMatchesOptions) ([]models.Match, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.PhaseId != "" {
		filter["phase_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PhaseId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "phase_id"
	}
	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	for cur.Next(ctx) {
		var match models.Match
		if err := cur.Decode(&match); err != nil {
			return nil, 0, err
		}
		matches = append(matches, match)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return matches, totalRecords, nil
}

func ProgramMatch(Time time.Time, Place string, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if !Time.IsZero() {
		updateDataMap["date"] = Time
	}
	if len(Place) > 0 {
		updateDataMap["place"] = Place
	}
	updateDataMap["status"] = models.Programmed

	return repositories.Update(match_collection, updateDataMap, Id)
}

func StartMatch(match models.Match, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["players_home"] = match.PlayersHome
	updateDataMap["players_away"] = match.PlayersAway
	updateDataMap["coachs_home"] = match.CoachsHome
	updateDataMap["coachs_away"] = match.CoachsAway
	updateDataMap["referees"] = match.Referees
	updateDataMap["scorekeeper"] = match.Scorekeeper
	updateDataMap["timekeeper"] = match.Timekeeper
	updateDataMap["status"] = models.FirstHalf

	return repositories.Update(match_collection, updateDataMap, Id)
}

func StartSecondHalf(Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["status"] = models.SecondHalf

	return repositories.Update(match_collection, updateDataMap, Id)
}

func EndMatch(Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["status"] = models.Ended

	return repositories.Update(match_collection, updateDataMap, Id)
}
