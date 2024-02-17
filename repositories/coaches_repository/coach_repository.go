package coaches_repository

import (
	"context"
	"strings"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	coach_collection = "coaches"
)

func CreateCoach(association_id string, coach models.Coach) (string, bool, error) {
	return repositories.Create(coach_collection, association_id, &coach)
}

func GetCoach(ID string) (models.Coach, bool, error) {
	var coach models.Coach
	_, err := repositories.GetById(coach_collection, ID, &coach)
	if err != nil {
		return models.Coach{}, false, err
	}

	return coach, true, nil
}

type GetCoachesOptions struct {
	Name          string
	Surname       string
	Dni           string
	Gender        string
	OnlyEnabled   bool
	TeamId        string
	AssociationId string
	Page          int
	PageSize      int
	SortOrder     int
}

func GetCoaches(filterOptions GetCoachesOptions) ([]models.Coach, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(coach_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		names := strings.Split(filterOptions.Name, " ")
		nameSurnameFilter := bson.A{}
		for _, name := range names {
			nameSurnameFilter = append(nameSurnameFilter, bson.M{"$or": []bson.M{
				{"personal_data.name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}},
				{"personal_data.surname": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}},
			}})
		}
		filter["$or"] = nameSurnameFilter
	}
	if filterOptions.Surname != "" {
		filter["personal_data.surname"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Surname, Options: "i"}}
	}
	if filterOptions.Dni != "" {
		filter["personal_data.dni"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Dni, Options: "i"}}
	}
	if filterOptions.Gender != "" {
		filter["personal_data.gender"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Gender, Options: "i"}}
	}
	if filterOptions.OnlyEnabled {
		filter["personal_data.disabled"] = false
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
		{Key: "personal_data.disabled", Value: sortOrder},
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

	var coaches []models.Coach
	for cur.Next(ctx) {
		var coach models.Coach
		if err := cur.Decode(&coach); err != nil {
			return nil, 0, err
		}
		coaches = append(coaches, coach)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return coaches, totalRecords, nil
}

func UpdateCoach(coach models.Coach, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(coach.Name) > 0 {
		updateDataMap["personal_data.name"] = coach.Name
	}
	if len(coach.Surname) > 0 {
		updateDataMap["personal_data.surname"] = coach.Surname
	}
	if !coach.DateOfBirth.IsZero() {
		updateDataMap["personal_data.date_of_birth"] = coach.DateOfBirth
	}
	if len(coach.Dni) > 0 {
		updateDataMap["personal_data.dni"] = coach.Dni
	}
	if len(coach.Dni) > 0 {
		updateDataMap["personal_data.gender"] = coach.Gender
	}
	if len(coach.PhoneNumber) > 0 {
		updateDataMap["personal_data.phone_number"] = coach.PhoneNumber
	}
	updateDataMap["personal_data.disabled"] = coach.Disabled

	if len(coach.TeamId) > 0 {
		updateDataMap["team_id"] = coach.TeamId
	}

	return repositories.Update(coach_collection, updateDataMap, ID)
}

func UpdateAvatar(coach models.Coach, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(coach.Avatar) > 0 {
		updateDataMap["personal_data.avatar"] = coach.Avatar
	}

	return repositories.Update(coach_collection, updateDataMap, ID)
}

func DeleteCoach(ID string) (bool, error) {
	return repositories.Delete(coach_collection, ID)
}

func GetCoachByDni(associationId, dni string) (models.Coach, bool, string) {
	condition := bson.M{"personal_data.dni": dni}

	var coach models.Coach
	_, err := repositories.FindOne(coach_collection, associationId, condition, &coach)
	id := coach.Id.Hex()
	if err != nil {
		return coach, false, id
	}
	return coach, true, id
}
