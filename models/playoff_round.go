package models

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ThirtyTwoFinals = "thirty_two_finals"
	SixteenFinals   = "sixteen_finals"
	EightFinals     = "eight_finals"
	QuarterFinals   = "quarter_finals"
	SemiFinal       = "semi_final"
	Final           = "final"
)

type PlayoffRound struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Round          string             `bson:"round" json:"round"`
	TeamsQuantity  int                `bson:"teams_quantity" json:"teams_quantity"`
	PlayoffPhaseId string             `bson:"playoff_phase_id" json:"playoff_phase_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
	AssociationId  string `bson:"association_id" json:"association_id"`
}

func (playoffRound *PlayoffRound) SetAssociationId(associationId string) {
	playoffRound.AssociationId = associationId
}

func (playoffRound *PlayoffRound) SetCreatedDate() {
	playoffRound.CreatedDate = time.Now()
}

func (playoffRound *PlayoffRound) SetModifiedDate() {
	playoffRound.ModifiedDate = time.Now()
}

func (playoffRound *PlayoffRound) SetId(id primitive.ObjectID) {
	playoffRound.Id = id
}

func CreatePlayoffRounds(playoffPhase PlayoffPhase) ([]PlayoffRound, []PlayoffRoundKey) {
	rounds, keys := createPlayoffRoundsRecursive(playoffPhase, playoffPhase.Teams, nil, nil)
	return rounds, keys
}

func createPlayoffRoundsRecursive(playoffPhase PlayoffPhase, teams []TournamentTeamId, rounds []PlayoffRound, keys []PlayoffRoundKey) ([]PlayoffRound, []PlayoffRoundKey) {
	if len(teams) <= 1 {
		return rounds, keys
	}

	round := PlayoffRound{
		Id:             primitive.NewObjectID(),
		Round:          GetRoundFromTeamsCount(len(teams)),
		TeamsQuantity:  len(teams),
		PlayoffPhaseId: playoffPhase.Id.Hex(),
	}

	roundKeys := make([]PlayoffRoundKey, len(teams)/2)
	for i := 0; i < len(teams)/2; i++ {
		keyNumber := i + 1
		key := PlayoffRoundKey{
			Id:             primitive.NewObjectID(),
			KeyNumber:      strconv.Itoa(keyNumber),
			PlayoffRoundId: round.Id.Hex(),
		}
		roundKeys[i] = key
	}
	keys = append(keys, roundKeys...)

	rounds = append(rounds, round)

	halfTeamsCount := len(teams) / 2
	return createPlayoffRoundsRecursive(playoffPhase, teams[:halfTeamsCount], rounds, keys)
}

func GetRoundFromTeamsCount(teamsCount int) string {
	switch {
	case teamsCount <= 2:
		return Final
	case teamsCount <= 4:
		return SemiFinal
	case teamsCount <= 8:
		return QuarterFinals
	case teamsCount <= 16:
		return EightFinals
	case teamsCount <= 32:
		return SixteenFinals
	default:
		return ThirtyTwoFinals
	}
}
