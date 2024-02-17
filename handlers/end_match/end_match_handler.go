package end_match

import (
	"github.com/nahuelojea/handballscore/handlers/end_match/domain"
	"github.com/nahuelojea/handballscore/models"
)

type EndMatchHandler interface {
	execute(*domain.EndMatch)
	setNext(EndMatchHandler)
}

func EndMatchHandlerExecute(match models.Match, leaguePhase models.LeaguePhase) {
	updateChampionHandler := &UpdateChampionHandler{}

	generateNewPhaseHandler := &GenerateNewPhaseHandler{}
	generateNewPhaseHandler.setNext(updateChampionHandler)

	updateTeamsScoreHandler := &UpdateTeamsScoreHandler{}
	updateTeamsScoreHandler.setNext(generateNewPhaseHandler)

	//endMatch := domain.EndMatch

	//updateTeamsScoreHandler.execute(&endMatch)
}
