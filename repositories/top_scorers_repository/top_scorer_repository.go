package top_scorers_repository

import (
	"context"
	"math"
	"time"

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
	// 1. Context con timeout para evitar bloqueos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := db.MongoClient.Database(db.DatabaseName).Collection("match_players_view")

	// Pre-cálculo de skip y limit
	skip := int64((filterOptions.Page - 1) * filterOptions.PageSize)
	limit := int64(filterOptions.PageSize)

	// 2. Pipeline optimizado
	pipeline := make([]bson.M, 0, 8)

	// 2.1. Filtrar por asociación y goles > 0 (usa índices)
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"association_id": filterOptions.AssociationId,
		"goals.total":    bson.M{"$gt": 0},
	}})

	// 2.2. Lookup solo de partidos relevantes (match_id y categoría)
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from": "matches",
		"let":  bson.M{"mid": "$match_id"},
		"pipeline": []bson.M{{
			"$match": bson.M{"$expr": bson.M{"$and": []bson.M{
				{"$eq": []interface{}{"$_id", bson.M{"$toObjectId": "$$mid"}}},
				{"$eq": []interface{}{"$tournament_category_id", filterOptions.TournamentCategoryId}},
				{"$not": bson.M{"$in": []interface{}{"$status", []string{models.Created, models.Programmed}}}},
			}}},
		}},
		"as": "match_info",
	}})

	// 2.3. Desenrollar y descartar no match
	pipeline = append(pipeline, bson.M{"$unwind": bson.M{
		"path":                       "$match_info",
		"preserveNullAndEmptyArrays": false,
	}})

	// 2.4. Agrupar goles y partidos por jugador
	pipeline = append(pipeline, bson.M{"$group": bson.M{
		"_id":            "$player_id",
		"total_goals":    bson.M{"$sum": "$goals.total"},
		"total_matches":  bson.M{"$addToSet": "$match_id"},
		"player_name":    bson.M{"$first": "$player_name"},
		"player_surname": bson.M{"$first": "$player_surname"},
		"player_avatar":  bson.M{"$first": "$player_avatar"},
		"team_name":      bson.M{"$first": "$team_name"},
		"team_avatar":    bson.M{"$first": "$team_avatar"},
	}})

	// 2.5. Filtrado opcional por nombre/apellido
	if filterOptions.Name != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$or": []bson.M{
			{"player_name": bson.M{"$regex": filterOptions.Name, "$options": "i"}},
			{"player_surname": bson.M{"$regex": filterOptions.Name, "$options": "i"}},
		}}})
	}

	// 2.6. Calcular total_matches y average
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"total_matches": bson.M{"$size": "$total_matches"},
		"average": bson.M{"$cond": []interface{}{
			bson.M{"$eq": []interface{}{bson.M{"$size": "$total_matches"}, 0}},
			0,
			bson.M{"$divide": []interface{}{"$total_goals", bson.M{"$size": "$total_matches"}}},
		}},
	}})

	// 2.7. Facet para datos y recuento en un solo pipeline
	pipeline = append(pipeline, bson.M{"$facet": bson.M{
		"data": []bson.M{
			{"$sort": bson.M{"total_goals": -1, "_id": 1}},
			{"$skip": skip},
			{"$limit": limit},
		},
		"total": []bson.M{
			{"$count": "totalRecords"},
		},
	}})

	// 2.8. Desenrollar total y proyectar resultado final
	pipeline = append(pipeline,
		bson.M{"$unwind": "$total"},
		bson.M{"$project": bson.M{
			"data":         1,
			"totalRecords": "$total.totalRecords",
		}},
	)

	// 3. Ejecutar agregación con allowDiskUse y batchSize
	opts := options.Aggregate().
		SetAllowDiskUse(true).
		SetBatchSize(100)

	cur, err := coll.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	// 4. Leer resultados
	var res struct {
		Data         []models.TopScorer `bson:"data"`
		TotalRecords int64              `bson:"totalRecords"`
	}
	if cur.Next(ctx) {
		if err := cur.Decode(&res); err != nil {
			return nil, 0, 0, err
		}
	}
	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	// 5. Calcular totalPages
	totalPages := int(math.Ceil(float64(res.TotalRecords) / float64(filterOptions.PageSize)))

	return res.Data, res.TotalRecords, totalPages, nil
}
