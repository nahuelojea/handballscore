package dto

type GetMatchPlayersRequest struct {
	TournamentTeamId TournamentTeamIdRequest `json:"team"`
	MatchId          string                  `json:"match_id"`
	PlayerId         string                  `json:"player_id"`
	HasBlueCard      bool                    `json:"has_blue_card"`
	Number           int                     `json:"number"`
}
