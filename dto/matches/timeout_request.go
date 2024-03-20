package dto

type TimeoutRequest struct {
	TournamentTeamId TournamentTeamIdRequest `json:"tournament_team_id"`
	Add              bool                    `json:"add"`
	Time             string                  `json:"time"`
}

type TournamentTeamIdRequest struct {
	TeamId  string `json:"team_id"`
	Variant string `json:"variant"`
}
