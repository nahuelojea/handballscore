package players_repository

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
	player_collection = "players"
)

func CreatePlayer(association_id string, player models.Player) (string, bool, error) {
	return repositories.Create(player_collection, association_id, &player)
}

func GetPlayer(ID string) (models.Player, bool, error) {
	var player models.Player
	_, err := repositories.GetById(player_collection, ID, &player)
	if err != nil {
		return models.Player{}, false, err
	}

	return player, true, nil
}

type GetPlayersOptions struct {
	Name          string
	Surname       string
	Dni           string
	Gender        string
	TeamId        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetPlayersFilteredAndPaginated(filterOptions GetPlayersOptions) ([]models.Player, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(player_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		filter["personal_data.name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
	}
	if filterOptions.Surname != "" {
		filter["personal_data.surname"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Surname, Options: "i"}}
	}
	if filterOptions.Dni != "" {
		filter["personal_data.dni"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Dni, Options: "i"}}
	}
	if filterOptions.Gender != "" {
		filter["gender"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Gender, Options: "i"}}
	}
	if filterOptions.TeamId != "" {
		filter["team_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.TeamId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	sortFields := bson.D{
		{Key: "status_data.disabled", Value: sortOrder},
		{Key: "personal_data.surname", Value: sortOrder},
		{Key: "personal_data.name", Value: sortOrder},
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

	var players []models.Player
	for cur.Next(ctx) {
		var player models.Player
		if err := cur.Decode(&player); err != nil {
			return nil, 0, err
		}
		players = append(players, player)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return players, totalRecords, nil
}

func UpdatePlayer(player models.Player, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(player.Name) > 0 {
		updateDataMap["personal_data.name"] = player.Name
	}
	if len(player.Surname) > 0 {
		updateDataMap["personal_data.surname"] = player.Surname
	}
	if len(player.Avatar) > 0 {
		updateDataMap["personal_data.avatar"] = player.Avatar
	}
	if !player.DateOfBirth.IsZero() {
		updateDataMap["personal_data.date_of_birth"] = player.DateOfBirth
	}
	if len(player.PhoneNumber) > 0 {
		updateDataMap["personal_data.phone_number"] = player.PhoneNumber
	}
	if len(player.AffiliateNumber) > 0 {
		updateDataMap["affiliate_number"] = player.AffiliateNumber
	}
	if len(player.Gender) > 0 {
		updateDataMap["gender"] = player.Gender
	}
	if !player.DateOfBirth.IsZero() {
		updateDataMap["expiration_insurance"] = player.ExpirationInsurance
	}
	if len(player.TeamId) > 0 {
		updateDataMap["team_id"] = player.TeamId
	}

	return repositories.Update(player_collection, updateDataMap, ID)
}

func DisablePlayer(ID string) (bool, error) {
	return repositories.Disable(player_collection, ID)
}

func GetPlayerByDni(dni string) (models.Player, bool, string) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(player_collection)

	condition := bson.M{"personal_data.dni": dni}

	var result models.Player

	err := collection.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()
	if err != nil {
		return result, false, id
	}
	return result, true, id
}