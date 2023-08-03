package referees_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	referee_collection = "referees"
)

func CreateReferee(referee models.Referee) (string, bool, error) {
	return repositories.Create(referee_collection, &referee)
}

func GetReferee(ID string) (models.Referee, error) {
	entity, err := repositories.GetById(referee_collection, ID)
	if err != nil {
		return models.Referee{}, err
	}

	referee, ok := entity.(*models.Referee)
	if !ok {
		return models.Referee{}, fmt.Errorf("Could not convert to Referee type")
	}

	return *referee, nil
}

type GetRefereesOptions struct {
	Name          string
	Surname       string
	Dni           string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetRefereesFilteredAndPaginated(filterOptions GetRefereesOptions) ([]models.Referee, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

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

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "personal_data.surname"
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

	var referees []models.Referee
	for cur.Next(ctx) {
		var referee models.Referee
		if err := cur.Decode(&referee); err != nil {
			return nil, 0, err
		}
		referees = append(referees, referee)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return referees, totalRecords, nil
}

func UpdateReferee(referee models.Referee, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(referee.Name) > 0 {
		updateDataMap["personal_data.name"] = referee.Name
	}
	if len(referee.Surname) > 0 {
		updateDataMap["personal_data.surname"] = referee.Surname
	}
	if len(referee.Avatar) > 0 {
		updateDataMap["personal_data.avatar"] = referee.Avatar
	}
	if !referee.DateOfBirth.IsZero() {
		updateDataMap["personal_data.date_of_birth"] = referee.DateOfBirth
	}
	if len(referee.Dni) > 0 {
		updateDataMap["personal_data.dni"] = referee.Dni
	}
	if len(referee.PhoneNumber) > 0 {
		updateDataMap["personal_data.phone_number"] = referee.PhoneNumber
	}
	updateDataMap["status_data.modified_date"] = time.Now()

	return repositories.Update(referee_collection, updateDataMap, ID)
}

func DisableReferee(ID string) (bool, error) {
	return repositories.Disable(referee_collection, ID)
}

func GetRefereeByDni(dni string) (models.Referee, bool, string) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(referee_collection)

	condition := bson.M{"personal_data.dni": dni}

	var result models.Referee

	err := collection.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()
	if err != nil {
		return result, false, id
	}
	return result, true, id
}
