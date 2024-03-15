package top_scorers_repository

import (
	"context"
	"log"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	match_player_collection = "match_players"
	player_collection       = "players"
	team_collection         = "teams"
)

func GetTopScorers(associationID string, limit int) ([]models.TopScorer, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	matchPlayerCollection := db.Collection(match_player_collection)
	teamCollection := db.Collection(team_collection)

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.M{"association_id": associationID}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":         "$player_id",
			"player_name": bson.M{"$first": "$player_name"},
			"avatar":      bson.M{"$first": "$avatar"},
			"team_id":     bson.M{"$first": "$team_id"},
			"total_goals": bson.M{"$sum": bson.M{"$add": []string{"$goals.first_half", "$goals.second_half"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"total_goals": -1}}},
		bson.D{{Key: "$limit", Value: limit}},
	}

	cur, err := matchPlayerCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)

	var topScorers []models.TopScorer
	for cur.Next(ctx) {
		var ts models.TopScorer
		err := cur.Decode(&ts)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Get player's team name
		var team models.Team
		err = teamCollection.FindOne(ctx, bson.M{"_id": ts.TeamId}).Decode(&team)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		ts.TeamName = team.Name

		topScorers = append(topScorers, ts)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return topScorers, nil
}
