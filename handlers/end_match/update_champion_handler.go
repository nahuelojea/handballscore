package end_match

import "github.com/nahuelojea/handballscore/handlers/end_match/domain"

type UpdateChampionHandler struct {
	next EndMatchHandler
}

func (updateChampionHandler *UpdateChampionHandler) execute(*domain.EndMatch) {

}

func (updateChampionHandler *UpdateChampionHandler) setNext(next EndMatchHandler) {
	updateChampionHandler.next = next
}
