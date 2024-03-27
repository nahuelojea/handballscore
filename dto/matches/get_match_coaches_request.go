package dto

type GetMatchCoachesRequest struct {
	TournamentTeamId TournamentTeamIdRequest `json:"team"`
	MatchId          string                  `json:"match_id"`
	CoachId          string                  `json:"coach_id"`
}
