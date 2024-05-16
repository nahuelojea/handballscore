package dto

type TimeoutRequest struct {
	TournamentTeamId TournamentTeamIdRequest `json:"team"`
	Add              bool                    `json:"add"`
	Time             string                  `json:"time"`
}

type TournamentTeamIdRequest struct {
	TeamId  string `json:"id"`
	Variant string `json:"variant"`
}
