package models

type TopScorer struct {
	PlayerId   string `bson:"_id"`
	PlayerName string `bson:"player_name"`
	Avatar     string `bson:"avatar"`
	TeamId     string `bson:"team_id"`
	TeamName   string `bson:"team_name"`
	Goals      int    `bson:"total_goals"`
}
