package dto

import "time"

type MatchResponse struct {
	MatchId   string            `bson:"match_id,omitempty" json:"match_id"`
	Date      time.Time         `bson:"date" json:"date"`
	TeamHome  MatchTeamResponse `bson:"team_home" json:"team_home"`
	TeamAway  MatchTeamResponse `bson:"team_away" json:"team_away"`
	Referees  []string          `bson:"referees" json:"referees"`
	Place     string            `bson:"place,omitempty" json:"place,omitempty"`
	PlaceId   string            `bson:"place_id,omitempty" json:"place_id,omitempty"`
	Status    string            `bson:"status" json:"status"`
	GoalsHome int               `bson:"goals_home" json:"goals_home"`
	GoalsAway int               `bson:"goals_away" json:"goals_away"`
	PlayoffRound string         `bson:"playoff_round" json:"playoff_round"`
}

type MatchTeamResponse struct {
	TeamId string `bson:"team_id" json:"team_id"`
	Variant string `bson:"variant" json:"variant"`
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
}
