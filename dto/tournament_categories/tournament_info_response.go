package dto

type TournamentInfoResponse struct {
	LeaguePhaseInfo  LeaguePhaseInfoResponse  `bson:"league_phase_info" json:"league_phase_info"`
	PlayoffPhaseInfo PlayoffPhaseInfoResponse `bson:"playoff_phase_info" json:"playoff_phase_info"`
}

type PlayoffPhaseInfoResponse struct {
	PlayoffKeys []PlayoffKeyResponse `bson:"playoff_keys" json:"playoff_keys"`
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

type PlayoffKeyResponse struct {
	Id               string                   `bson:"id" json:"id"`
	Name             string                   `bson:"name" json:"name"`
	NextPlayoffKeyId string                   `bson:"next_playoff_key_id" json:"next_playoff_key_id"`
	State            string                   `bson:"state" json:"state"`
	PlayoffKeyTeams  []PlayoffKeyTeamResponse `bson:"playoff_key_teams" json:"playoff_key_teams"`
}

type PlayoffKeyTeamResponse struct {
	Id string `bson:"id" json:"id"`
	TeamInfoResponse
	Result string `bson:"result" json:"result"`
	Status string `bson:"status" json:"status"`
}
