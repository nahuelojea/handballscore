package top_scorers_repository

import (
	"github.com/nahuelojea/handballscore/models"
)

const (
	match_player_collection = "match_players"
	player_collection       = "players"
	team_collection         = "teams"
)

type GetTopScorersOptions struct {
	TournamentCategoryId string
	AssociationId        string
	Page                 int
	PageSize             int
	SortField            string
	SortOrder            int
}

func GetTopScorers(filterOptions GetTopScorersOptions) ([]models.TopScorer, int64, int, error) {
	/*ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(tournament_category_collection)

	filter := bson.M{
		"association_id": filterOptions.AssociationId,
	}

	if filterOptions.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Name, Options: "i"}}
	}
	if filterOptions.CategoryId != "" {
		filter["category_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.CategoryId, Options: "i"}}
	}
	if filterOptions.TournamentId != "" {
		filter["tournament_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.TournamentId, Options: "i"}}
	}
	if filterOptions.Status != "" {
		filter["status"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.Status, Options: "i"}}
	}
	if filterOptions.ChampionId != "" {
		filter["champion_id"] = bson.M{"$regex": primitive.Regex{Pattern: filterOptions.ChampionId, Options: "i"}}
	}

	page := filterOptions.Page
	pageSize := filterOptions.PageSize

	sortOrder := 1
	if filterOptions.SortOrder == 1 {
		sortOrder = 1
	}

	sortFields := bson.D{
		{Key: "start_date", Value: sortOrder},
		{Key: "end_date", Value: sortOrder},
		{Key: "status", Value: sortOrder},
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetSort(sortFields)

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var tournaments []models.TournamentCategory
	for cur.Next(ctx) {
		var tournament models.TournamentCategory
		if err := cur.Decode(&tournament); err != nil {
			return nil, 0, 0, err
		}
		tournaments = append(tournaments, tournament)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return tournaments, totalRecords, totalPages, nil*/

	// Create sample top scorers for testing
	topScorers := []models.TopScorer{
		{PlayerName: "John Doe",
			PlayerAvatar: "https://handballscore.s3.amazonaws.com/avatars/players/6545542c5cf49946f89ef98a.jpg",
			TeamName:     "Handball Norte",
			TeamAvatar:   "https://handballscore.s3.amazonaws.com/avatars/teams/6537cda15f3e10c95cc1726f.jpg",
			Goals:        40,
			Matches:      5,
			Average:      8},
		{PlayerName: "Jane Smith",
			PlayerAvatar: "https://handballscore.s3.amazonaws.com/avatars/players/6545542c5cf49946f89ef98a.jpg",
			TeamName:     "Handball Sur",
			TeamAvatar:   "https://handballscore.s3.amazonaws.com/avatars/teams/6537cda15f3e10c95cc1726f.jpg",
			Goals:        35,
			Matches:      6,
			Average:      5},
		{PlayerName: "Mike Johnson",
			PlayerAvatar: "https://handballscore.s3.amazonaws.com/avatars/players/6545542c5cf49946f89ef98a.jpg",
			TeamName:     "Handball Este",
			TeamAvatar:   "https://handballscore.s3.amazonaws.com/avatars/teams/6537cda15f3e10c95cc1726f.jpg",
			Goals:        30,
			Matches:      4,
			Average:      7},
		{PlayerName: "Sarah Davis",
			PlayerAvatar: "https://handballscore.s3.amazonaws.com/avatars/players/6545542c5cf49946f89ef98a.jpg",
			TeamName:     "Handball Oeste",
			TeamAvatar:   "https://handballscore.s3.amazonaws.com/avatars/teams/6537cda15f3e10c95cc1726f.jpg",
			Goals:        25,
			Matches:      3,
			Average:      8},
		{PlayerName: "David Wilson",
			PlayerAvatar: "https://handballscore.s3.amazonaws.com/avatars/players/6545542c5cf49946f89ef98a.jpg",
			TeamName:     "Handball Centro",
			TeamAvatar:   "https://handballscore.s3.amazonaws.com/avatars/teams/6537cda15f3e10c95cc1726f.jpg",
			Goals:        20,
			Matches:      2,
			Average:      10},
		{PlayerName: "Emily Thompson",
			PlayerAvatar: "https://handballscore.s3.amazonaws.com/avatars/players/6545542c5cf49946f89ef98a.jpg",
			TeamName:     "Handball Sur",
			TeamAvatar:   "https://handballscore.s3.amazonaws.com/avatars/teams/6537cda15f3e10c95cc1726f.jpg",
			Goals:        15,
			Matches:      1,
			Average:      15},
	}

	return topScorers, 6, 1, nil
}
