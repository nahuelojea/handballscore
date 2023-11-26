package domain

import "github.com/nahuelojea/handballscore/models"

type EndMatch struct {
	match            models.Match
	leaguePhase      models.LeaguePhase
	updateTeamsScore StepStatus
	generateNewPhase StepStatus
	updateChampion   StepStatus
}

type StepStatus struct {
	isDone bool
	status string
}
