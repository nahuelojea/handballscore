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
	MatchId  string           `bson:"match_id" json:"match_id"`
	PlayerId string           `bson:"player_id" json:"player_id"`
	Team     MatchTeamRequest `bson:"team" json:"team"`
	Number   string           `bson:"number" json:"number"`
}

type MatchCoachRequest struct {
	MatchId string           `bson:"match_id" json:"match_id"`
	CoachId string           `bson:"coach_id" json:"coach_id"`
	Team    MatchTeamRequest `bson:"team" json:"team"`
}

type MatchTeamRequest struct {
	Id      string `bson:"id" json:"id"`
	Variant string `bson:"variant" json:"variant"`
}
