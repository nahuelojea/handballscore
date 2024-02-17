package match_players_repository

import (
	"context"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	match_player_collection = "match_players"
)

func CreateMatchPlayer(association_id string, matchPlayer models.MatchPlayer) (string, bool, error) {
	return repositories.Create(match_player_collection, association_id, &matchPlayer)
}

func CreateMatchPlayers(association_id string, matchPlayers []models.MatchPlayer) ([]string, bool, error) {
	entities := make([]models.Entity, len(matchPlayers))
	for i, v := range matchPlayers {
		matchPlayer := v
		entities[i] = models.Entity(&matchPlayer)
	}
	return repositories.CreateMultiple(match_player_collection, association_id, entities)
}

func GetMatchPlayer(id string) (models.MatchPlayer, bool, error) {
	var matchPlayer models.MatchPlayer
	_, err := repositories.GetById(match_player_collection, id, &matchPlayer)
	if err != nil {
		return models.MatchPlayer{}, false, err
	}
	return matchPlayer, true, nil
}

type GetMatchPlayerOptions struct {
	MatchId       string
	TeamId        string
	PlayerId      string
	Number        string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchPlayers(filterOptions GetMatchPlayerOptions) ([]models.MatchPlayer, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_player_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.MatchId != "" {
		filter["match_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.MatchId, Options: "i"}}
	}
	if filterOptions.TeamId != "" {
		filter["team_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.TeamId, Options: "i"}}
	}
	if filterOptions.PlayerId != "" {
		filter["player_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PlayerId, Options: "i"}}
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
		{Key: "player_id", Value: sortOrder},
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(sortFields)

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var matchPlayers []models.MatchPlayer
	for cur.Next(ctx) {
		var matchPlayer models.MatchPlayer
		if err := cur.Decode(&matchPlayer); err != nil {
			return nil, 0, err
		}
		matchPlayers = append(matchPlayers, matchPlayer)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return matchPlayers, totalRecords, nil
}

func UpdateMatchPlayer(matchPlayer models.MatchPlayer, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	if len(matchPlayer.Number) > 0 {
		updateDataMap["number"] = matchPlayer.Number
	}

	return repositories.Update(match_player_collection, updateDataMap, id)
}

func UpdateGoals(matchPlayer models.MatchPlayer, status string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	if status == models.FirstHalf {
		updateDataMap["goals.first_half"] = matchPlayer.Goals.FirstHalf
	} else {
		updateDataMap["goals.second_half"] = matchPlayer.Goals.SecondHalf
	}

	return repositories.Update(match_player_collection, updateDataMap, matchPlayer.Id.Hex())
}

func UpdateExclusions(matchPlayer models.MatchPlayer) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.exclusions"] = matchPlayer.Exclusions

	return repositories.Update(match_player_collection, updateDataMap, matchPlayer.Id.Hex())
}

func UpdateYellowCard(matchPlayer models.MatchPlayer) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.yellow_card"] = matchPlayer.YellowCard

	return repositories.Update(match_player_collection, updateDataMap, matchPlayer.Id.Hex())
}

func UpdateRedCard(matchPlayer models.MatchPlayer) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.red_card"] = matchPlayer.RedCard

	return repositories.Update(match_player_collection, updateDataMap, matchPlayer.Id.Hex())
}

func UpdateBlueCard(matchPlayer models.MatchPlayer) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["sanctions.blue_card"] = matchPlayer.BlueCard
	updateDataMap["sanctions.report"] = matchPlayer.Report

	return repositories.Update(match_player_collection, updateDataMap, matchPlayer.Id.Hex())
}

func DeleteMatchPlayer(id string) (bool, error) {
	return repositories.Delete(match_player_collection, id)
}
