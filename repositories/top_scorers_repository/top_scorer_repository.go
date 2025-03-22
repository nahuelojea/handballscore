package top_scorers_repository

import (
	"context"
	"math"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetTopScorersOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Name                 string
	Page                 int
	PageSize             int
}

func GetTopScorers(filterOptions GetTopScorersOptions) ([]models.TopScorer, int64, int, error) {
	ctx := context.TODO()
	database := db.MongoClient.Database(db.DatabaseName)
	collection := database.Collection("match_players_view")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"association_id": filterOptions.AssociationId,
				"goals.total":    bson.M{"$gt": 0},
			},
		},
		{
			"$lookup": bson.M{
				"from": "matches",
				"let":  bson.M{"match_id": "$match_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$_id", bson.M{"$toObjectId": "$$match_id"}}},
									{"$eq": []interface{}{"$tournament_category_id", filterOptions.TournamentCategoryId}},
									{"$not": bson.M{"$in": []interface{}{"$status", []string{models.Created, models.Programmed}}}},
								},
							},
						},
					},
					{"$limit": 1},
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
			"$group": bson.M{
				"_id":            "$player_id",
				"total_goals":    bson.M{"$sum": "$goals.total"},
				"total_matches":  bson.M{"$addToSet": "$match_id"},
				"player_name":    bson.M{"$first": "$player_name"},
				"player_surname": bson.M{"$first": "$player_surname"},
				"player_avatar":  bson.M{"$first": "$player_avatar"},
				"team_name":      bson.M{"$first": "$team_name"},
				"team_avatar":    bson.M{"$first": "$team_avatar"},
			},
		},
	}

	if filterOptions.Name != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"$or": []bson.M{
					{"player_name": bson.M{"$regex": filterOptions.Name, "$options": "i"}},
					{"player_surname": bson.M{"$regex": filterOptions.Name, "$options": "i"}},
				},
			},
		})
	}

	facetPipeline := []bson.M{
		{
			"$addFields": bson.M{
				"total_matches": bson.M{"$size": "$total_matches"},
				"average": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []interface{}{bson.M{"$size": "$total_matches"}, 0}},
						0,
						bson.M{"$divide": []interface{}{"$total_goals", bson.M{"$size": "$total_matches"}}},
					},
				},
			},
		},
		{
			"$sort": bson.M{
				"total_goals": -1,
				"_id":         1,
			},
		},
		{
			"$facet": bson.M{
				"data": []bson.M{
					{"$skip": (filterOptions.Page - 1) * filterOptions.PageSize},
					{"$limit": filterOptions.PageSize},
				},
				"total": []bson.M{
					{"$count": "totalRecords"},
				},
			},
		},
		{
			"$unwind": "$total",
		},
	}

	pipeline = append(pipeline, facetPipeline...)

	opts := options.Aggregate().SetAllowDiskUse(true)
	cur, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var result struct {
		Data  []models.TopScorer `bson:"data"`
		Total struct {
			TotalRecords int64 `bson:"totalRecords"`
		} `bson:"total"`
	}

	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			return nil, 0, 0, err
		}
	}

	totalRecords := int64(0)
	if len(result.Data) > 0 {
		totalRecords = result.Total.TotalRecords
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(filterOptions.PageSize)))

	return result.Data, totalRecords, totalPages, nil
}
