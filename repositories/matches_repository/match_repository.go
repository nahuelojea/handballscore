package matches_repository

import (
	"context"
	"math"
	"time"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	match_collection  = "matches"
	match_header_view = "match_headers_view"
)

func CreateMatches(associationID string, matches []models.Match) ([]string, bool, error) {
	entities := make([]models.Entity, len(matches))
	for i, v := range matches {
		match := v
		entities[i] = models.Entity(&match)
	}

	return repositories.CreateMultiple(match_collection, associationID, entities)
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

func GetMatchHeaderView(ID string) (models.MatchHeaderView, bool, error) {
	var matchHeader models.MatchHeaderView
	_, err := repositories.GetById(match_header_view, ID, &matchHeader)
	if err != nil {
		return models.MatchHeaderView{}, false, err
	}

	return matchHeader, true, nil
}

type GetMatchesOptions struct {
	LeaguePhaseWeekId  string
	PlayoffRoundKeyIds []string
	Date               time.Time
	AssociationId      string
	Page               int
	PageSize           int
	SortField          string
	SortOrder          int
}

func GetMatches(filterOptions GetMatchesOptions) ([]models.Match, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.LeaguePhaseWeekId != "" {
		filter["league_phase_week_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.LeaguePhaseWeekId, Options: "i"}}
	}

	if len(filterOptions.PlayoffRoundKeyIds) > 0 {
		filter["playoff_round_key_id"] = bson.M{"$in": filterOptions.PlayoffRoundKeyIds}
	}

	if !filterOptions.Date.IsZero() {
		startDate := time.Date(filterOptions.Date.Year(), filterOptions.Date.Month(), filterOptions.Date.Day(), 0, 0, 0, 0, filterOptions.Date.Location())
		endDate := startDate.AddDate(0, 0, 1)
		filter["date"] = bson.M{"$gte": startDate, "$lt": endDate}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "date"
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

	var matches []models.Match
	for cur.Next(ctx) {
		var match models.Match
		if err := cur.Decode(&match); err != nil {
			return nil, 0, 0, err
		}
		matches = append(matches, match)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return matches, totalRecords, totalPages, nil
}

func GetMatchHeaders(filterOptions GetMatchesOptions) ([]models.MatchHeaderView, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_header_view)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if !filterOptions.Date.IsZero() {
		startDate := time.Date(filterOptions.Date.Year(), filterOptions.Date.Month(), filterOptions.Date.Day(), 0, 0, 0, 0, filterOptions.Date.Location())
		endDate := startDate.AddDate(0, 0, 1)
		filter["date"] = bson.M{"$gte": startDate, "$lt": endDate}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortField := filterOptions.SortField
	if sortField == "" {
		sortField = "date"
	}
	sortOrder := 1
	if filterOptions.SortOrder == -1 {
		sortOrder = -1
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(bson.D{
		{Key: "place", Value: sortOrder},
		{Key: "date", Value: sortOrder},
	})

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var matchViews []models.MatchHeaderView
	for cur.Next(ctx) {
		var matchView models.MatchHeaderView
		if err := cur.Decode(&matchView); err != nil {
			return nil, 0, 0, err
		}
		matchViews = append(matchViews, matchView)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return matchViews, totalRecords, totalPages, nil
}

func ProgramMatch(time time.Time, place, id string) (bool, error) {
	updateDataMap := make(map[string]interface{})
	if !time.IsZero() {
		updateDataMap["date"] = time
	}
	if len(place) > 0 {
		updateDataMap["place"] = place
	}

	updateDataMap["status"] = models.Programmed

	return repositories.Update(match_collection, updateDataMap, id)
}

func UpdateReferees(match models.Match, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["referees"] = match.Referees

	return repositories.Update(match_collection, updateDataMap, Id)
}

func StartMatch(match models.Match, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["scorekeeper"] = match.Scorekeeper
	updateDataMap["timekeeper"] = match.Timekeeper
	updateDataMap["status"] = models.FirstHalf

	return repositories.Update(match_collection, updateDataMap, Id)
}

func StartSecondHalf(Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["status"] = models.SecondHalf

	return repositories.Update(match_collection, updateDataMap, Id)
}

func EndMatch(id, comments string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["status"] = models.Finished
	updateDataMap["comments"] = comments

	return repositories.Update(match_collection, updateDataMap, id)
}

func SuspendMatch(id, comments string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["status"] = models.Suspended
	updateDataMap["comments"] = comments

	return repositories.Update(match_collection, updateDataMap, id)
}

func UpdateGoals(match models.Match, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["goals_home"] = match.GoalsHome
	updateDataMap["goals_away"] = match.GoalsAway

	return repositories.Update(match_collection, updateDataMap, Id)
}

func UpdateTimeouts(match models.Match, Id string) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["timeouts_home"] = match.TimeoutsHome
	updateDataMap["timeouts_away"] = match.TimeoutsAway

	return repositories.Update(match_collection, updateDataMap, Id)
}

func GetPendingMatchesByLeaguePhaseId(leaguePhaseId string) ([]models.Match, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from": "league_phase_weeks",
				"let": bson.M{
					"league_phase_week_id_str": "$league_phase_week_id",
				},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": bson.A{
									"$_id",
									bson.M{"$toObjectId": "$$league_phase_week_id_str"},
								},
							},
						},
					},
				},
				"as": "league_phase_week_info",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$league_phase_week_info",
				"preserveNullAndEmptyArrays": false,
			},
		},
		bson.M{
			"$match": bson.M{
				"league_phase_week_info.league_phase_id": leaguePhaseId,
				"status":                                 bson.M{"$nin": bson.A{models.Finished, models.Suspended}},
			},
		},
	}

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	for cur.Next(ctx) {
		var match models.Match
		if err := cur.Decode(&match); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func GetPendingMatchesByPlayoffRoundKeyId(playoffRoundKeyId string) ([]models.Match, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	filter := bson.M{
		"playoff_round_key_id": playoffRoundKeyId,
		"status": bson.M{
			"$nin": bson.A{models.Finished, models.Suspended},
		},
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	for cur.Next(ctx) {
		var match models.Match
		if err := cur.Decode(&match); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func UpdateHomeTeam(id string, team models.TournamentTeamId) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["team_home"] = team

	return repositories.Update(match_collection, updateDataMap, id)
}

func UpdateAwayTeam(id string, team models.TournamentTeamId) (bool, error) {
	updateDataMap := make(map[string]interface{})

	updateDataMap["team_away"] = team

	return repositories.Update(match_collection, updateDataMap, id)
}

func GetLastEndedMatchByTeam(team models.TournamentTeamId, tournamentCategoryId string) (models.Match, bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(match_collection)

	filter := bson.M{
		"$and": bson.A{
			bson.M{"$or": bson.A{
				bson.M{"team_home.team_id": team.TeamId, "team_home.variant": team.Variant, "status": models.Finished},
				bson.M{"team_away.team_id": team.TeamId, "team_away.variant": team.Variant, "status": models.Finished},
			}},
			bson.M{"tournament_category_id": tournamentCategoryId},
		},
	}

	sort := bson.D{{Key: "date", Value: -1}}

	findOptions := options.FindOne()
	findOptions.SetSort(sort)

	var match models.Match
	err := collection.FindOne(ctx, filter, findOptions).Decode(&match)
	if err != nil {
		return match, false, err
	}

	return match, true, nil
}
