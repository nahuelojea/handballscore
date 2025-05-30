package teams_repository

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
	team_collection = "teams"
)

type GetTeamsOptions struct {
	Name          string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateTeam(association_id string, team models.Team) (string, bool, error) {
	return repositories.Create(team_collection, association_id, &team)
}

func GetTeam(ID string) (models.Team, bool, error) {
	var team models.Team
	_, err := repositories.GetById(team_collection, ID, &team)
	if err != nil {
		return models.Team{}, false, err
	}

	return team, true, nil
}

func GetTeams(filterOptions GetTeamsOptions) ([]models.Team, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(team_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
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
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var teams []models.Team
	for cur.Next(ctx) {
		var team models.Team
		if err := cur.Decode(&team); err != nil {
			return nil, 0, 0, err
		}
		teams = append(teams, team)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return teams, totalRecords, totalPages, nil
}

func UpdateTeam(team models.Team, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(team.Name) > 0 {
		updateDataMap["name"] = team.Name
	}
	if len(team.Address) > 0 {
		updateDataMap["address"] = team.Address
	}
	if len(team.PhoneNumber) > 0 {
		updateDataMap["phone_number"] = team.PhoneNumber
	}
	if !team.DateOfFoundation.IsZero() {
		updateDataMap["date_of_foundation"] = team.DateOfFoundation
	}
	if len(team.Initials) > 0 {
		updateDataMap["initials"] = team.Initials
	}
	if len(team.MainColor) > 0 {
		updateDataMap["main_color"] = team.MainColor
	}
	if len(team.SecondaryColor) > 0 {
		updateDataMap["secondary_color"] = team.SecondaryColor
	}
	if len(team.Email) > 0 {
		updateDataMap["email"] = team.Email
	}

	return repositories.Update(team_collection, updateDataMap, ID)
}

func UpdateAvatar(team models.Team, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(team.Avatar) > 0 {
		updateDataMap["avatar"] = team.Avatar
	}

	return repositories.Update(team_collection, updateDataMap, ID)
}

func DeleteTeam(ID string) (bool, error) {
	return repositories.Delete(team_collection, ID)
}
