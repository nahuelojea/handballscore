package player_sanctions_repository

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
	player_sanction_collection = "player_sanctions"
)

type GetPlayerSanctionsOptions struct {
	PlayerId       string
	AssociationId  string
	SanctionStatus string
	Page           int
	PageSize       int
	SortField      string
	SortOrder      int
}

func CreatePlayerSanction(associationId string, playerSanction models.PlayerSanction) (string, bool, error) {
	return repositories.Create(player_sanction_collection, associationId, &playerSanction)
}

func GetPlayerSanction(id string) (models.PlayerSanction, bool, error) {
	var playerSanction models.PlayerSanction
	_, err := repositories.GetById(player_sanction_collection, id, &playerSanction)
	if err != nil {
		return models.PlayerSanction{}, false, err
	}

	return playerSanction, true, nil
}

func UpdatePlayerSanction(playerSanction models.PlayerSanction, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if !playerSanction.IssueDate.IsZero() {
		updateDataMap["issue_date"] = playerSanction.IssueDate
	}
	if !playerSanction.EndDate.IsZero() {
		updateDataMap["end_date"] = playerSanction.EndDate
	}
	if len(playerSanction.Description) > 0 {
		updateDataMap["description"] = playerSanction.Description
	}
	if playerSanction.MatchSuspensions > 0 {
		updateDataMap["match_suspensions"] = playerSanction.MatchSuspensions
	}
	if len(playerSanction.SanctionStatus) > 0 {
		updateDataMap["sanction_status"] = playerSanction.SanctionStatus
	}

	return repositories.Update(player_sanction_collection, updateDataMap, id)
}

func AddServedMatch(matchId, id string) (bool, error) {
	playerSanction, _, err := GetPlayerSanction(id)
	if err != nil {
		return false, err
	}

	updateDataMap := make(map[string]interface{})

	playerSanction.ServedMatches = append(playerSanction.ServedMatches, matchId)
	updateDataMap["served_matches"] = playerSanction.ServedMatches

	if len(playerSanction.ServedMatches) == playerSanction.MatchSuspensions {
		updateDataMap["sanction_status"] = models.Completed
	}

	return repositories.Update(player_sanction_collection, updateDataMap, id)
}

func GetPlayerSanctions(filterOptions GetPlayerSanctionsOptions) ([]models.PlayerSanction, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(player_sanction_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if len(filterOptions.PlayerId) > 0 {
		filter["player_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PlayerId, Options: "i"}}
	}

	if len(filterOptions.SanctionStatus) > 0 {
		filter["sanction_status"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.SanctionStatus, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "issue_date"
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
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var playerSanctions []models.PlayerSanction
	for cur.Next(ctx) {
		var playerSanction models.PlayerSanction
		if err := cur.Decode(&playerSanction); err != nil {
			return nil, 0, 0, err
		}
		playerSanctions = append(playerSanctions, playerSanction)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return playerSanctions, totalRecords, totalPages, nil
}
