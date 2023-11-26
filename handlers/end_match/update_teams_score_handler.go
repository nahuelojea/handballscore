package end_match

import "github.com/nahuelojea/handballscore/handlers/end_match/domain"

type UpdateTeamsScoreHandler struct {
	next EndMatchHandler
}

func (updateTeamsScoreHandler *UpdateTeamsScoreHandler) execute(*domain.EndMatch) {

}

func (updateTeamsScoreHandler *UpdateTeamsScoreHandler) setNext(next EndMatchHandler) {
	updateTeamsScoreHandler.next = next
}
