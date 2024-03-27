package dto

type GetMatchPlayersRequest struct {
	TournamentTeamId TournamentTeamIdRequest `json:"team"`
	MatchId          string                  `json:"match_id"`
	PlayerId         string                  `json:"player_id"`
	Number           string                  `json:"number"`
}
