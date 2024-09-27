package top_scorers_repository

import (
	"context"
	"math"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
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
			},
		},
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
									{"$eq": []interface{}{"$tournament_category_id", filterOptions.TournamentCategoryId}},
									{"$not": bson.M{
										"$in": []interface{}{"$status", []string{models.Created, models.Programmed}},
									}},
								},
							},
						},
					},
					{
						"$project": bson.M{
							"date":      1,
							"team_home": 1,
							"team_away": 1,
							"place":     1,
							"status":    1,
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
			"$match": bson.M{
				"goals.total": bson.M{"$gt": 0}, // Filtrar jugadores con goles
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
				"association_id": bson.M{"$first": "$association_id"},
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

	pipeline = append(pipeline,
		bson.M{
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
		bson.M{
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
		bson.M{
			"$sort": bson.M{
				"total_goals": -1,
			},
		},
		bson.M{
			"$skip": int64((filterOptions.Page - 1) * filterOptions.PageSize),
		},
		bson.M{
			"$limit": int64(filterOptions.PageSize),
		},
	)

	cur, err := collection.Aggregate(ctx, pipeline)
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

	countPipeline := append(pipeline[:len(pipeline)-3], bson.M{
		"$count": "totalRecords",
	})

	var countResult struct {
		TotalRecords int64 `bson:"totalRecords"`
	}
	countCur, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, 0, err
	}
	defer countCur.Close(ctx)

	if countCur.Next(ctx) {
		if err := countCur.Decode(&countResult); err != nil {
			return nil, 0, 0, err
		}
	}

	totalPages := int(math.Ceil(float64(countResult.TotalRecords) / float64(filterOptions.PageSize)))

	return topScorers, countResult.TotalRecords, totalPages, nil
}
