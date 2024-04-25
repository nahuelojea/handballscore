package models

type TopScorer struct {
	PlayerName    string  `bson:"player_name"`
	PlayerAvatar  string  `bson:"player_avatar"`
	TeamName      string  `bson:"team_name"`
	TeamAvatar    string  `bson:"team_avatar"`
	Goals         int     `bson:"total_goals"`
	Matches       int     `bson:"total_matches"`
	Average       float64 `bson:"average"`
	AssociationId string  `bson:"association_id"`
}
