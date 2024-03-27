package match_coaches_repository

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
	match_coach_collection = "match_coaches"
)

func CreateMatchCoach(association_id string, matchCoach models.MatchCoach) (string, bool, error) {
	return repositories.Create(match_coach_collection, association_id, &matchCoach)
}

func CreateMatchCoaches(association_id string, matchCoach []models.MatchCoach) ([]string, bool, error) {
	entities := make([]models.Entity, len(matchCoach))
	for i, v := range matchCoach {
		matchCoach := v
		entities[i] = models.Entity(&matchCoach)
	}
	return repositories.CreateMultiple(match_coach_collection, association_id, entities)
}

func GetMatchCoach(id string) (models.MatchCoach, bool, error) {
	var matchCoach models.MatchCoach
	_, err := repositories.GetById(match_coach_collection, id, &matchCoach)
	if err != nil {
		return models.MatchCoach{}, false, err
	}
	return matchCoach, true, nil
}

type GetMatchCoachOptions struct {
	MatchId       string
	Team          models.TournamentTeamId
	CoachId       string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchCoaches(filterOptions GetMatchCoachOptions) ([]models.MatchCoach, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_coach_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.MatchId != "" {
		filter["match_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.MatchId, Options: "i"}}
	}
	if filterOptions.Team.TeamId != "" {
		filter["team.team_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Team.TeamId, Options: "i"}}
	}
	if filterOptions.Team.Variant != "" {
		filter["team.variant"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Team.Variant, Options: "i"}}
	}
	if filterOptions.CoachId != "" {
		filter["coach_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.CoachId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "coach_id", Value: sortOrder},
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

	var matchCoaches []models.MatchCoach
	for cur.Next(ctx) {
		var matchCoach models.MatchCoach
		if err := cur.Decode(&matchCoach); err != nil {
			return nil, 0, 0, err
		}
		matchCoaches = append(matchCoaches, matchCoach)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return matchCoaches, totalRecords, totalPages, nil
}

func UpdateExclusions(matchCoach models.MatchCoach) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.exclusions"] = matchCoach.Exclusions

	return repositories.Update(match_coach_collection, updateDataMap, matchCoach.Id.Hex())
}

func UpdateYellowCard(matchCoach models.MatchCoach) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.yellow_card"] = matchCoach.YellowCard

	return repositories.Update(match_coach_collection, updateDataMap, matchCoach.Id.Hex())
}

func UpdateRedCard(matchCoach models.MatchCoach) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.red_card"] = matchCoach.RedCard

	return repositories.Update(match_coach_collection, updateDataMap, matchCoach.Id.Hex())
}

func UpdateBlueCard(matchCoach models.MatchCoach) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.blue_card"] = matchCoach.BlueCard
	updateDataMap["sanctions.report"] = matchCoach.Report

	return repositories.Update(match_coach_collection, updateDataMap, matchCoach.Id.Hex())
}

func DeleteMatchCoach(id string) (bool, error) {
	return repositories.Delete(match_coach_collection, id)
}
