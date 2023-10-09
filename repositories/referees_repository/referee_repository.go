package referees_repository

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
	referee_collection = "referees"
)

func CreateReferee(association_id string, referee models.Referee) (string, bool, error) {
	return repositories.Create(referee_collection, association_id, &referee)
}

func GetReferee(ID string) (models.Referee, error) {
	var referee models.Referee
	_, err := repositories.GetById(referee_collection, ID, &referee)
	if err != nil {
		return models.Referee{}, err
	}

	return referee, nil
}

type GetRefereesOptions struct {
	Name          string
	Surname       string
	Dni           string
	Gender        string
	OnlyEnabled   bool
	AssociationId string
	Page          int
	PageSize      int
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
	if filterOptions.Gender != "" {
		filter["personal_data.gender"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Gender, Options: "i"}}
	}
	if filterOptions.OnlyEnabled {
		filter["personal_data.disabled"] = false
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
	if !referee.DateOfBirth.IsZero() {
		updateDataMap["personal_data.date_of_birth"] = referee.DateOfBirth
	}
	if len(referee.Dni) > 0 {
		updateDataMap["personal_data.dni"] = referee.Dni
	}
	if len(referee.Dni) > 0 {
		updateDataMap["personal_data.gender"] = referee.Gender
	}
	if len(referee.PhoneNumber) > 0 {
		updateDataMap["personal_data.phone_number"] = referee.PhoneNumber
	}
	updateDataMap["personal_data.disabled"] = referee.Disabled

	return repositories.Update(referee_collection, updateDataMap, ID)
}

func UpdateAvatar(referee models.Referee, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(referee.Avatar) > 0 {
		updateDataMap["personal_data.avatar"] = referee.Avatar
	}

	return repositories.Update(referee_collection, updateDataMap, ID)
}

func DeleteReferee(ID string) (bool, error) {
	return repositories.Delete(referee_collection, ID)
}

func GetRefereeByDni(associationId, dni string) (models.Referee, bool, string) {
	condition := bson.M{"personal_data.dni": dni}

	var referee models.Referee
	_, err := repositories.FindOne(referee_collection, associationId, condition, &referee)
	id := referee.Id.Hex()
	if err != nil {
		return referee, false, id
	}
	return referee, true, id
}
