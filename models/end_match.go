package models

const (
	League_Current_Phase  = "league"
	Playoff_Current_Phase = "playoff"
)

type EndMatch struct {
	Match                      Match
	CurrentLeaguePhase         CurrentLeaguePhase
	CurrentPlayoffPhase        CurrentPlayoffPhase
	CurrentTournamentCategory  TournamentCategory
	CurrentPhase               string
	UpdateTeamsScore           StepStatus
	GenerateNewPlayoffRoundKey StepStatus
	GenerateNewPhase           StepStatus
	UpdateChampion             StepStatus
	LoadDataNextMatches        StepStatus
	UpdatePlayersSanctions     StepStatus
}

type StepStatus struct {
	IsDone bool
	Status string
}

type CurrentLeaguePhase struct {
	LeaguePhase     LeaguePhase
	LeaguePhaseWeek LeaguePhaseWeek
}

type CurrentPlayoffPhase struct {
	PlayoffPhase    PlayoffPhase
	PlayoffRound    PlayoffRound
	PlayoffRoundKey PlayoffRoundKey
}
