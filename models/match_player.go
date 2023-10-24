package models

type MatchPlayer struct {
	PlayerId  string `bson:"player_id" json:"player_id"`
	MatchId   string `bson:"match_id" json:"match_id"`
	Number    string `bson:"number" json:"number"`
	Goals     `bson:"goals" json:"goals"`
	Sanctions `bson:"sanctions" json:"sanctions"`
}

type Goals struct {
	FirstHalf  int `bson:"first_half" json:"first_half"`
	SecondHalf int `bson:"second_half" json:"second_half"`
}
