package dto

type WeeksAndRoundsResponse struct {
	Description       string `bson:"description" json:"description"`
	LeaguePhaseWeekId string `bson:"league_phase_week_id" json:"league_phase_week_id"`
	PlayoffRoundId    string `bson:"playoff_round_id" json:"playoff_round_id"`
}
