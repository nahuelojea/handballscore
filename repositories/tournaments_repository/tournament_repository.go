package tournaments_repository

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
	tournament_collection = "tournaments"
)

func CreateTournament(association_id string, tournament models.TournamentCategory) (string, bool, error) {
	return repositories.Create(tournament_collection, association_id, &tournament)
}

func GetTournament(ID string) (models.TournamentCategory, bool, error) {
	var tournament models.TournamentCategory
	_, err := repositories.GetById(tournament_collection, ID, &tournament)
	if err != nil {
		return models.TournamentCategory{}, false, err
	}

	return tournament, true, nil
}

type GetTournamentsOptions struct {
	Name          string
	CategoryId    string
	Status        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetTournamentsFilteredAndPaginated(filterOptions GetTournamentsOptions) ([]models.TournamentCategory, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(tournament_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
	}
	if filterOptions.CategoryId != "" {
		filter["category_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.CategoryId, Options: "i"}}
	}
	if filterOptions.CategoryId != "" {
		filter["status"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Status, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "name"
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

	var tournaments []models.TournamentCategory
	for cur.Next(ctx) {
		var tournament models.TournamentCategory
		if err := cur.Decode(&tournament); err != nil {
			return nil, 0, err
		}
		tournaments = append(tournaments, tournament)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return tournaments, totalRecords, nil
}

func UpdateTournament(tournament models.TournamentCategory, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(tournament.Name) > 0 {
		updateDataMap["name"] = tournament.Name
	}
	if len(tournament.Champion) > 0 {
		updateDataMap["champion"] = tournament.Champion
	}
	if len(tournament.Status) > 0 {
		updateDataMap["status"] = tournament.Status
	}

	return repositories.Update(tournament_collection, updateDataMap, ID)
}

func DeleteTournament(ID string) (bool, error) {
	return repositories.Delete(tournament_collection, ID)
}
