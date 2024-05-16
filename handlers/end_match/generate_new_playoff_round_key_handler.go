package end_match

import (
	"github.com/nahuelojea/handballscore/models"
)

type GenerateNewPlayoffRoundKeyHandler struct {
	BaseEndMatchHandler
}

func (c *GenerateNewPlayoffRoundKeyHandler) HandleEndMatch(endMatch *models.EndMatch) {

	if c.GetNext() != nil {
		c.GetNext().HandleEndMatch(endMatch)
	}
}
