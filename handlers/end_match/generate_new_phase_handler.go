package end_match

import "github.com/nahuelojea/handballscore/handlers/end_match/domain"

type GenerateNewPhaseHandler struct {
	next EndMatchHandler
}

func (generateNewPhaseHandler *GenerateNewPhaseHandler) execute(*domain.EndMatch) {

}

func (generateNewPhaseHandler *GenerateNewPhaseHandler) setNext(next EndMatchHandler) {
	generateNewPhaseHandler.next = next
}
