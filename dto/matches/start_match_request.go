package dto

type StartMatchRequest struct {
	PlayersHome []MatchPlayerRequest `bson:"players_home" json:"players_home"`
	CoachsHome  []string             `bson:"coachs_home" json:"coachs_home"`
	PlayersAway []MatchPlayerRequest `bson:"players_away" json:"players_away"`
	CoachsAway  []string             `bson:"coachs_away" json:"coachs_away"`
	Referees    []string             `bson:"referees" json:"referees"`
	Scorekeeper string               `bson:"scorekeeper" json:"scorekeeper"`
	Timekeeper  string               `bson:"timekeeper" json:"timekeeper"`
}

type MatchPlayerRequest struct {
	PlayerId string `bson:"id" json:"id"`
	Number   string `bson:"number" json:"number"`
}
