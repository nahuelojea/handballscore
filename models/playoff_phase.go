package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayoffPhase struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams                []TournamentTeamId `bson:"teams" json:"teams"`
	Config               PlayoffPhaseConfig `bson:"config" json:"config"`
	TournamentCategoryId string             `bson:"tournament_category_id" json:"tournament_category_id"`
	Status_Data          `bson:"status_data" json:"status_data"`
	AssociationId        string `bson:"association_id" json:"association_id"`
}

type PlayoffPhaseConfig struct {
	HomeAndAway      bool `bson:"home_and_away" json:"home_and_away"`
	SingleMatchFinal bool `bson:"single_match_final" json:"single_match_final"`
	RandomOrder      bool `bson:"random_order" json:"random_order"`
	ClassifiedNumber int  `bson:"classified_number" json:"classified_number"`
}

func (playoffPhase *PlayoffPhase) SetAssociationId(associationId string) {
	playoffPhase.AssociationId = associationId
}

func (playoffPhase *PlayoffPhase) SetCreatedDate() {
	playoffPhase.CreatedDate = time.Now()
}

func (playoffPhase *PlayoffPhase) SetModifiedDate() {
	playoffPhase.ModifiedDate = time.Now()
}

func (playoffPhase *PlayoffPhase) SetId(id primitive.ObjectID) {
	playoffPhase.Id = id
}

func CreateRoundMatches(playoffPhase PlayoffPhase, roundKeys []PlayoffRoundKey) []Match {
	matches := []Match{}

	for i := 0; i < len(roundKeys); i++ {
		teamA := roundKeys[i].Teams[0]
		teamB := roundKeys[i].Teams[1]

		match := GeneratePlayoffMatch(playoffPhase.TournamentCategoryId, roundKeys[i].Id.Hex(), teamA, teamB)
		matches = append(matches, match)

		if playoffPhase.Config.HomeAndAway {
			matchReturn := GeneratePlayoffMatch(playoffPhase.TournamentCategoryId, roundKeys[i].Id.Hex(), teamB, teamA)
			matches = append(matches, matchReturn)
		}
	}

	return matches
}
