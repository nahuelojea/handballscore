package dto

import "time"

type MatchJourneyResponse struct {
	MatchId   string            `bson:"match_id,omitempty" json:"match_id"`
	Date      time.Time         `bson:"date" json:"date"`
	TeamHome  MatchTeamResponse `bson:"team_home" json:"team_home"`
	TeamAway  MatchTeamResponse `bson:"team_away" json:"team_away"`
	Place     string            `bson:"place" json:"place"`
	Status    string            `bson:"status" json:"status"`
	GoalsHome int               `bson:"goals_home" json:"goals_home"`
	GoalsAway int               `bson:"goals_away" json:"goals_away"`
}

type MatchTeamResponse struct {
	TeamId string `bson:"team_id" json:"team_id"`
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
}
