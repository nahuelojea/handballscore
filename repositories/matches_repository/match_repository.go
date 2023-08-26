package matches_repository

import (
	"context"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	match_collection = "matches"
)

func CreateMatches(associationID string, matches []models.Match) ([]string, bool, error) {
	var createdIDs []string

	for _, match := range matches {
		match.SetCreatedDate()
		match.SetModifiedDate()
		match.SetDisabled(false)
		match.SetAssociationId(associationID)

		id, created, err := repositories.Create(match_collection, associationID, &match)
		if err != nil {
			return createdIDs, false, err
		}
		if created {
			createdIDs = append(createdIDs, id)
		}
	}

	return createdIDs, true, nil
}

func CreateMatch(association_id string, match models.Match) (string, bool, error) {
	return repositories.Create(match_collection, association_id, &match)
}

func GetMatch(ID string) (models.Match, bool, error) {
	var match models.Match
	_, err := repositories.GetById(match_collection, ID, &match)
	if err != nil {
		return models.Match{}, false, err
	}

	return match, true, nil
}

type GetMatchesOptions struct {
	PhaseId       string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func GetMatchesFilteredAndPaginated(filterOptions GetMatchesOptions) ([]models.Match, int64, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.PhaseId != "" {
		filter["phase_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.PhaseId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "phase_id"
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

	var matches []models.Match
	for cur.Next(ctx) {
		var match models.Match
		if err := cur.Decode(&match); err != nil {
			return nil, 0, err
		}
		matches = append(matches, match)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return matches, totalRecords, nil
}

func ProgramMatch(Time time.Time, Place string, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if !Time.IsZero() {
		updateDataMap["date"] = Time
	}
	if len(Place) > 0 {
		updateDataMap["place"] = Place
	}
	updateDataMap["status"] = models.Programmed

	return repositories.Update(match_collection, updateDataMap, Id)
}

func StartMatch(startMatchRequest dto.StartMatchRequest, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["players_local"] = startMatchRequest.PlayersLocal
	updateDataMap["players_visiting"] = startMatchRequest.PlayersVisiting
	updateDataMap["coachs_local"] = startMatchRequest.CoachsLocal
	updateDataMap["coachs_visiting"] = startMatchRequest.CoachsVisiting
	updateDataMap["referees"] = startMatchRequest.Referees
	updateDataMap["scorekeeper"] = startMatchRequest.Scorekeeper
	updateDataMap["timekeeper"] = startMatchRequest.Timekeeper
	updateDataMap["status"] = models.FirstHalf

	return repositories.Update(match_collection, updateDataMap, Id)
}

/*func MatchGoal(matchGoalRequest dto.MatchGoalRequest, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["players_local"] = startMatchRequest.PlayersLocal
	updateDataMap["players_visiting"] = startMatchRequest.PlayersVisiting
	updateDataMap["coachs_local"] = startMatchRequest.CoachsLocal
	updateDataMap["coachs_visiting"] = startMatchRequest.CoachsVisiting
	updateDataMap["referees"] = startMatchRequest.Referees
	updateDataMap["scorekeeper"] = startMatchRequest.Scorekeeper
	updateDataMap["timekeeper"] = startMatchRequest.Timekeeper
	updateDataMap["status"] = models.FirstHalf

	return repositories.Update(match_collection, updateDataMap, Id)
}*/

/*func UpdateMatch(match models.Match, ID string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if len(match.Name) > 0 {
		updateDataMap["name"] = match.Name
	}
	if len(match.Surname) > 0 {
		updateDataMap["personal_data.surname"] = match.Surname
	}
	if len(match.Avatar) > 0 {
		updateDataMap["personal_data.avatar"] = match.Avatar
	}
	if !match.DateOfBirth.IsZero() {
		updateDataMap["personal_data.date_of_birth"] = match.DateOfBirth
	}
	if len(match.Dni) > 0 {
		updateDataMap["personal_data.dni"] = match.Dni
	}
	if len(match.PhoneNumber) > 0 {
		updateDataMap["personal_data.phone_number"] = match.PhoneNumber
	}
	if len(match.AffiliateNumber) > 0 {
		updateDataMap["affiliate_number"] = match.AffiliateNumber
	}
	if len(match.Gender) > 0 {
		updateDataMap["gender"] = match.Gender
	}
	if len(match.TeamId) > 0 {
		updateDataMap["team_id"] = match.TeamId
	}

	return repositories.Update(match_collection, updateDataMap, ID)
}

func DisableMatch(ID string) (bool, error) {
	return repositories.Disable(match_collection, ID)
}

func GetMatchByDni(dni string) (models.Match, bool, string) {
	ctx := context.TODO()

	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	condition := bson.M{"personal_data.dni": dni}

	var result models.Match

	err := collection.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()
	if err != nil {
		return result, false, id
	}
	return result, true, id
}*/
