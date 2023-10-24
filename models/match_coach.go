package models

type MatchCoach struct {
	CoachId   string `bson:"coach_id" json:"coach_id"`
	MatchId   string `bson:"match_id" json:"match_id"`
	Sanctions `bson:"sanctions" json:"sanctions"`
}
