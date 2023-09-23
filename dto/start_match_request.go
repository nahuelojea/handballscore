package dto

import (
	"github.com/nahuelojea/handballscore/models"
)

type StartMatchRequest struct {
	PlayersLocal    []models.MatchPlayer  `bson:"players_local" json:"players_local"`
	CoachsLocal     []models.MatchCoach   `bson:"coachs_local" json:"coachs_local"`
	PlayersVisiting []models.MatchPlayer  `bson:"players_visiting" json:"players_visiting"`
	CoachsVisiting  []models.MatchCoach   `bson:"coachs_visiting" json:"coachs_visiting"`
	Referees        []models.MatchReferee `bson:"referees" json:"referees"`
	Scorekeeper     string                `bson:"scorekeeper" json:"scorekeeper"`
	Timekeeper      string                `bson:"timekeeper" json:"timekeeper"`
}
