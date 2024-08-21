package top_scorers_repository

import (
	"context"
	"math"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	match_players_view = "match_players_view"
	top_scorers_view   = "top_scorers_view"
)

type GetTopScorersOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func CreateTopScorersView(ctx context.Context, db *mongo.Database, tournamentCategoryId string) error {
	viewName := "top_scorers_view"

	if err := db.RunCommand(context.Background(), bson.D{{Key: "drop", Value: viewName}}).Err(); err != nil {
		if err.Error() != "ns not found" {
			return err
		}
	}

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": "matches",
				"let":  bson.M{"match_id_str": "$match_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$_id", bson.M{"$toObjectId": "$$match_id_str"}}},
									{"$eq": []interface{}{"$tournament_category_id", tournamentCategoryId}},
									{"$not": bson.M{
										"$in": []interface{}{"$status", []string{models.Created, models.Programmed}},
									}},
								},
							},
						},
					},
				},
				"as": "match_info",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$match_info",
				"preserveNullAndEmptyArrays": false,
			},
		},
		{
			"$project": bson.M{
				"match_id":       "$match_id",
				"player_id":      "$player_id",
				"goals":          "$goals.total",
				"date":           "$match_info.date",
				"team_home":      "$match_info.team_home",
				"team_away":      "$match_info.team_away",
				"place":          "$match_info.place",
				"status":         "$match_info.status",
				"player_name":    "$player_name",
				"player_surname": "$player_surname",
				"player_avatar":  "$player_avatar",
				"team_name":      "$team_name",
				"team_avatar":    "$team_avatar",
				"association_id": "$association_id",
			},
		},
		{
			"$group": bson.M{
				"_id":            "$player_id",
				"total_goals":    bson.M{"$sum": "$goals"},
				"total_matches":  bson.M{"$addToSet": "$match_id"},
				"player_name":    bson.M{"$first": "$player_name"},
				"player_surname": bson.M{"$first": "$player_surname"},
				"player_avatar":  bson.M{"$first": "$player_avatar"},
				"team_name":      bson.M{"$first": "$team_name"},
				"team_avatar":    bson.M{"$first": "$team_avatar"},
				"association_id": bson.M{"$first": "$association_id"},
			},
		},
		{
			"$match": bson.M{
				"total_goals": bson.M{"$gt": 0},
			},
		},
		{
			"$addFields": bson.M{
				"total_matches": bson.M{"$size": "$total_matches"},
				"average": bson.M{
					"$divide": []interface{}{
						"$total_goals",
						bson.M{"$size": "$total_matches"},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"total_goals":    1,
				"average":        1,
				"total_matches":  1,
				"player_name":    1,
				"player_surname": 1,
				"player_avatar":  1,
				"team_name":      1,
				"team_avatar":    1,
				"association_id": 1,
			},
		},
		{
			"$sort": bson.M{
				"total_goals": -1,
			},
		},
	}

	err := db.CreateView(ctx, viewName, "match_players_view", pipeline)
	if err != nil {
		return err
	}

	return nil
}

func GetTopScorers(filterOptions GetTopScorersOptions) ([]models.TopScorer, int64, int, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)

	err := CreateTopScorersView(ctx, db, filterOptions.TournamentCategoryId)
	if err != nil {
		return nil, 0, 0, err
	}

	collection := db.Collection(top_scorers_view)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var topScorers []models.TopScorer
	for cur.Next(ctx) {
		var topScorer models.TopScorer
		if err := cur.Decode(&topScorer); err != nil {
			return nil, 0, 0, err
		}
		topScorers = append(topScorers, topScorer)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return topScorers, totalRecords, totalPages, nil
}
