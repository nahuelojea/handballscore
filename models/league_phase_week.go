package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaguePhaseWeek struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number        int                `bson:"number" json:"number"`
	Status_Data   `bson:"status_data" json:"status_data"`
	LeaguePhaseId string `bson:"league_phase_id" json:"league_phase_id"`
	AssociationId string `bson:"association_id" json:"association_id"`
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetAssociationId(associationId string) {
	leaguePhaseWeek.AssociationId = associationId
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetCreatedDate() {
	leaguePhaseWeek.CreatedDate = time.Now()
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetModifiedDate() {
	leaguePhaseWeek.ModifiedDate = time.Now()
}

func (leaguePhaseWeek *LeaguePhaseWeek) SetId(id primitive.ObjectID) {
	leaguePhaseWeek.Id = id
}

func GenerateLeaguePhaseWeeks(leaguePhase LeaguePhase) []LeaguePhaseWeek {
	totalTeams := len(leaguePhase.Teams)
	var weeks int

	if leaguePhase.HomeAndAway {
		weeks = totalTeams*2 - 2
	} else {
		weeks = totalTeams - 1
	}

	var leaguePhaseWeeks []LeaguePhaseWeek

	for i := 1; i <= weeks; i++ {
		leaguePhaseWeek := LeaguePhaseWeek{
			Number:        i,
			LeaguePhaseId: leaguePhase.Id.Hex(),
		}
		leaguePhaseWeeks = append(leaguePhaseWeeks, leaguePhaseWeek)
	}

	return leaguePhaseWeeks
}

func GenerateLeaguePhaseWeek(leaguePhase LeaguePhase, weekNumber int) LeaguePhaseWeek {
	return LeaguePhaseWeek{
		Number:        weekNumber,
		LeaguePhaseId: leaguePhase.Id.Hex(),
		AssociationId: leaguePhase.AssociationId,
	}
}
