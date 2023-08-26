package dto

type MatchGoalRequest struct {
	PlayerId string `bson:"coachs_local" json:"coachs_local"`
}
