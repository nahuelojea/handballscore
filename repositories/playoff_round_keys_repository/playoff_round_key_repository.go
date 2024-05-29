package playoff_round_keys_repository

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
	playoff_round_key_collection = "playoff_round_keys"
)

type GetPlayoffRoundKeysOptions struct {
	PlayoffRoundId string
	AssociationId  string
	Page           int
	PageSize       int
	SortField      string
	SortOrder      int
}

func CreatePlayoffRoundKey(association_id string, playoffRoundKey models.PlayoffRoundKey) (string, bool, error) {
	return repositories.Create(playoff_round_key_collection, association_id, &playoffRoundKey)
}

func CreatePlayoffRoundKeys(association_id string, playoffRoundKeys []models.PlayoffRoundKey) ([]string, bool, error) {
	entities := make([]models.Entity, len(playoffRoundKeys))
	for i, v := range playoffRoundKeys {
		playoffRoundKey := v
		entities[i] = models.Entity(&playoffRoundKey)
	}
	return repositories.CreateMultiple(playoff_round_key_collection, association_id, entities)
}

func GetPlayoffRoundKey(ID string) (models.PlayoffRoundKey, bool, error) {
	var playoffRoundKey models.PlayoffRoundKey
	_, err := repositories.GetById(playoff_round_key_collection, ID, &playoffRoundKey)
	if err != nil {
		return models.PlayoffRoundKey{}, false, err
	}

	return playoffRoundKey, true, nil
}

func GetPlayoffRoundKeys(filterOptions GetPlayoffRoundKeysOptions) ([]models.PlayoffRoundKey, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(playoff_round_key_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.PlayoffRoundId != "" {
		filter["playoff_round_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PlayoffRoundId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "playoff_round_id", Value: sortOrder},
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

	var playoffRoundKeys []models.PlayoffRoundKey
	for cur.Next(ctx) {
		var playoffRoundKey models.PlayoffRoundKey
		if err := cur.Decode(&playoffRoundKey); err != nil {
			return nil, 0, 0, err
		}
		playoffRoundKeys = append(playoffRoundKeys, playoffRoundKey)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return playoffRoundKeys, totalRecords, totalPages, nil
}

func UpdateTeamsRanking(playoffRoundKey models.PlayoffRoundKey, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["teams_ranking"] = playoffRoundKey.TeamsRanking

	return repositories.Update(playoff_round_key_collection, updateDataMap, id)
}

func DeletePlayoffRoundKey(ID string) (bool, error) {
	return repositories.Delete(playoff_round_key_collection, ID)
}
