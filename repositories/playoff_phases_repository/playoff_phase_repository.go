package playoff_phases_repository

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
	playoff_phase_collection = "playoff_phases"
)

type GetPlayoffPhasesOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func CreatePlayoffPhase(association_id string, playoffPhase models.PlayoffPhase) (string, bool, error) {
	return repositories.Create(playoff_phase_collection, association_id, &playoffPhase)
}

func GetPlayoffPhase(ID string) (models.PlayoffPhase, bool, error) {
	var playoffPhase models.PlayoffPhase
	_, err := repositories.GetById(playoff_phase_collection, ID, &playoffPhase)
	if err != nil {
		return models.PlayoffPhase{}, false, err
	}

	return playoffPhase, true, nil
}

func GetPlayoffPhases(filterOptions GetPlayoffPhasesOptions) ([]models.PlayoffPhase, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(playoff_phase_collection)

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

	var playoffPhases []models.PlayoffPhase
	for cur.Next(ctx) {
		var playoffPhase models.PlayoffPhase
		if err := cur.Decode(&playoffPhase); err != nil {
			return nil, 0, 0, err
		}
		playoffPhases = append(playoffPhases, playoffPhase)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return playoffPhases, totalRecords, totalPages, nil
}

func DeletePlayoffPhase(ID string) (bool, error) {
	return repositories.Delete(playoff_phase_collection, ID)
}
