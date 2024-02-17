package tournaments

import "time"

const (
	League_Format             = "league"
	Playoff_Format            = "playoff"
	League_And_Playoff_Format = "league_and_playoff"
)

type CreateTournamentCategoryRequest struct {
	CategoryId   string              `json:"category_id"`
	TournamentId string              `json:"tournament_id"`
	StartDate    time.Time           `json:"start_date"`
	Format       string              `json:"format"`
	LeaguePhase  LeaguePhaseRequest  `json:"league_phase"`
	PlayoffPhase PlayoffPhaseRequest `json:"playoff_phase"`
	Teams        []string            `json:"teams"`
}

type LeaguePhaseRequest struct {
	HomeAndAway      bool `json:"home_and_away"`
	ClassifiedNumber int  `json:"classified_number"`
}

type PlayoffPhaseRequest struct {
	HomeAndAway bool `json:"home_and_away"`
	RandomOrder bool `json:"random_order"`
}
