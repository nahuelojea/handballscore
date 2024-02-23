package playoff_rounds_repository

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
	playoff_round_collection = "playoff_rounds"
)

type GetPlayoffRoundsOptions struct {
	PlayoffPhaseId string
	AssociationId  string
	Page           int
	PageSize       int
	SortField      string
	SortOrder      int
}

func CreatePlayoffRound(association_id string, playoffRound models.PlayoffRound) (string, bool, error) {
	return repositories.Create(playoff_round_collection, association_id, &playoffRound)
}

func CreatePlayoffRounds(association_id string, playoffRounds []models.PlayoffRound) ([]string, bool, error) {
	entities := make([]models.Entity, len(playoffRounds))
	for i, v := range playoffRounds {
		playoffRound := v
		entities[i] = models.Entity(&playoffRound)
	}
	return repositories.CreateMultiple(playoff_round_collection, association_id, entities)
}

func GetPlayoffRound(ID string) (models.PlayoffRound, bool, error) {
	var playoffRound models.PlayoffRound
	_, err := repositories.GetById(playoff_round_collection, ID, &playoffRound)
	if err != nil {
		return models.PlayoffRound{}, false, err
	}

	return playoffRound, true, nil
}

func GetPlayoffRounds(filterOptions GetPlayoffRoundsOptions) ([]models.PlayoffRound, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(playoff_round_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.PlayoffPhaseId != "" {
		filter["playoff_phase_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PlayoffPhaseId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "playoff_phase_id", Value: sortOrder},
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

	var playoffRounds []models.PlayoffRound
	for cur.Next(ctx) {
		var playoffRound models.PlayoffRound
		if err := cur.Decode(&playoffRound); err != nil {
			return nil, 0, 0, err
		}
		playoffRounds = append(playoffRounds, playoffRound)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return playoffRounds, totalRecords, totalPages, nil
}

func DeletePlayoffRound(ID string) (bool, error) {
	return repositories.Delete(playoff_round_collection, ID)
}
