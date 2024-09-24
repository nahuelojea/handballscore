package match_players_repository

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
	match_player_collection = "match_players"
	match_player_view       = "match_players_view"
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
	Team          models.TournamentTeamId
	PlayerId      string
	Number        int
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetMatchPlayers(filterOptions GetMatchPlayerOptions) ([]models.MatchPlayerView, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_player_view)

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
	if filterOptions.PlayerId != "" {
		filter["player_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PlayerId, Options: "i"}}
	}
	if filterOptions.Number != 0 {
		filter["number"] = filterOptions.Number
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
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var matchPlayersView []models.MatchPlayerView
	for cur.Next(ctx) {
		var matchPlayerView models.MatchPlayerView
		if err := cur.Decode(&matchPlayerView); err != nil {
			return nil, 0, 0, err
		}
		matchPlayersView = append(matchPlayersView, matchPlayerView)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return matchPlayersView, totalRecords, totalPages, nil
}

func UpdateMatchPlayer(matchPlayer models.MatchPlayer, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["number"] = matchPlayer.Number

	return repositories.Update(match_player_collection, updateDataMap, id)
}

func UpdateGoals(matchPlayer models.MatchPlayer) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["goals.first_half"] = matchPlayer.Goals.FirstHalf
	updateDataMap["goals.second_half"] = matchPlayer.Goals.SecondHalf
	updateDataMap["goals.total"] = matchPlayer.Goals.Total

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
