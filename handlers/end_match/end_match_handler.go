package end_match

import (
	"github.com/nahuelojea/handballscore/models"
)

type EndMatchHandler interface {
	HandleEndMatch(*models.EndMatch)
}

type BaseEndMatchHandler struct {
	nextHandler EndMatchHandler
}

func (h *BaseEndMatchHandler) SetNext(next EndMatchHandler) {
	h.nextHandler = next
}

func (h *BaseEndMatchHandler) GetNext() EndMatchHandler {
	return h.nextHandler
}

func (h *BaseEndMatchHandler) HandleEndMatch(endMatch *models.EndMatch) {
	if h.nextHandler != nil {
		h.nextHandler.HandleEndMatch(endMatch)
	}
}
