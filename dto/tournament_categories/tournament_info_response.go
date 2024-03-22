package dto

type TournamentInfoResponse struct {
	LeaguePhaseInfo  LeaguePhaseInfoResponse  `bson:"league_phase_info" json:"league_phase_info"`
	PlayoffPhaseInfo PlayoffPhaseInfoResponse `bson:"playoff_phase_info" json:"playoff_phase_info"`
}

type PlayoffPhaseInfoResponse struct {
	PlayoffRounds []PlayoffRoundInfoResponse `bson:"playoff_rounds" json:"playoff_rounds"`
}

type PlayoffRoundInfoResponse struct {
	Round               string                        `bson:"round" json:"round"`
	PlayoffRoundKeyInfo []PlayoffRoundKeyInfoResponse `bson:"playoff_round_key_info" json:"playoff_round_key_info"`
}

type PlayoffRoundKeyInfoResponse struct {
	TeamsRanking []TeamScoreResponse `bson:"teams_ranking" json:"teams_ranking"`
}

type LeaguePhaseInfoResponse struct {
	TeamsRanking []TeamScoreResponse `bson:"teams_ranking" json:"teams_ranking"`
}

type TeamScoreResponse struct {
	TeamInfo        TeamInfoResponse `bson:"team_info" json:"team_info"`
	Position        int              `bson:"position" json:"position"`
	Classified      bool             `bson:"classified" json:"classified"`
	Points          int              `bson:"points" json:"points"`
	Matches         int              `bson:"matches" json:"matches"`
	Wins            int              `bson:"wins" json:"wins"`
	Draws           int              `bson:"draws" json:"draws"`
	Losses          int              `bson:"losses" json:"losses"`
	GoalsScored     int              `bson:"goals_scored" json:"goals_scored"`
	GoalsConceded   int              `bson:"goals_conceded" json:"goals_conceded"`
	GoalsDifference int              `bson:"goals_difference" json:"goals_difference"`
}

type TeamInfoResponse struct {
	TeamName   string `bson:"team_name" json:"team_name"`
	TeamAvatar string `bson:"team_avatar" json:"team_avatar"`
}
