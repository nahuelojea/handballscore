package tournaments_repository

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
	tournament_category_collection = "tournaments_categories"
	tournament_category_view       = "tournaments_categories_view"
)

func CreateTournamentCategory(association_id string, tournament models.TournamentCategory) (string, bool, error) {
	return repositories.Create(tournament_category_collection, association_id, &tournament)
}

func GetTournamentCategory(ID string) (models.TournamentCategory, bool, error) {
	var tournamentCategory models.TournamentCategory
	_, err := repositories.GetById(tournament_category_collection, ID, &tournamentCategory)
	if err != nil {
		return models.TournamentCategory{}, false, err
	}

	return tournamentCategory, true, nil
}

type GetTournamentsCategoryOptions struct {
	Name          string
	CategoryId    string
	TournamentId  string
	Status        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetTournamentsCategories(filterOptions GetTournamentsCategoryOptions) ([]models.TournamentCategoryView, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(tournament_category_view)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
	}
	if filterOptions.CategoryId != "" {
		filter["category_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.CategoryId, Options: "i"}}
	}
	if filterOptions.TournamentId != "" {
		filter["tournament_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.TournamentId, Options: "i"}}
	}
	if filterOptions.Status != "" {
		filter["status"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Status, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := -1
	if filterOptions.SortOrder == 1 {
		sortOrder = 1
	}

	sortFields := bson.D{
		{Key: "start_date", Value: sortOrder},
		{Key: "end_date", Value: sortOrder},
		{Key: "status", Value: sortOrder},
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

	var tournaments []models.TournamentCategoryView
	for cur.Next(ctx) {
		var tournament models.TournamentCategoryView
		if err := cur.Decode(&tournament); err != nil {
			return nil, 0, 0, err
		}
		tournaments = append(tournaments, tournament)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return tournaments, totalRecords, totalPages, nil
}

func UpdateTournamentCategory(tournament models.TournamentCategory, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(tournament.Name) > 0 {
		updateDataMap["name"] = tournament.Name
	}
	if len(tournament.Champion.TeamId) > 0 {
		updateDataMap["champion.team_id"] = tournament.Champion.TeamId
	}
	if len(tournament.Champion.Variant) > 0 {
		updateDataMap["champion.variant"] = tournament.Champion.Variant
	}
	if len(tournament.Status) > 0 {
		updateDataMap["status"] = tournament.Status
	}
	if !tournament.StartDate.IsZero() {
		updateDataMap["start_date"] = tournament.StartDate
	}
	if !tournament.EndDate.IsZero() {
		updateDataMap["end_date"] = tournament.EndDate
	}

	return repositories.Update(tournament_category_collection, updateDataMap, id)
}

func DeleteTournamentCategory(ID string) (bool, error) {
	return repositories.Delete(tournament_category_collection, ID)
}
