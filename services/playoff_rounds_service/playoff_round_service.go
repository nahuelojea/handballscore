package playoff_rounds_service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/playoff_rounds_repository"
	"github.com/nahuelojea/handballscore/services/matches_service"
	"github.com/nahuelojea/handballscore/services/playoff_round_keys_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePlayoffRound(association_id string, playoffRound models.PlayoffRound) (string, bool, error) {
	return playoff_rounds_repository.CreatePlayoffRound(association_id, playoffRound)
}

func CreatePlayoffRounds(association_id string, playoffRounds []models.PlayoffRound) ([]string, bool, error) {
	return playoff_rounds_repository.CreatePlayoffRounds(association_id, playoffRounds)
}

func CreateTournamentPlayoffRounds(association_id string, playoffPhase models.PlayoffPhase) (string, bool, error) {

	playoffRounds, playoffRoundKeys, matches := createRounds(playoffPhase)

	_, _, err := CreatePlayoffRounds(association_id, playoffRounds)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff rounds: %s", err.Error()))
	}

	_, _, err = playoff_round_keys_service.CreatePlayoffRoundKeys(association_id, playoffRoundKeys)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff round keys: %s", err.Error()))
	}

	_, _, err = matches_service.CreateMatches(association_id, matches)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("Error to create playoff matches: %s", err.Error()))
	}

	return "", true, nil
}

func createRounds(playoffPhase models.PlayoffPhase) ([]models.PlayoffRound, []models.PlayoffRoundKey, []models.Match) {
	if playoffPhase.Config.ClassifiedNumber == 0 {
		rounds, keys, matches := createPlayoffRoundsRecursive(playoffPhase, playoffPhase.Teams, len(playoffPhase.Teams), nil, nil, nil)
		return rounds, keys, matches
	} else {
		rounds, keys, matches := createPlayoffRoundsRecursive(playoffPhase, playoffPhase.Teams, playoffPhase.Config.ClassifiedNumber, nil, nil, nil)
		return rounds, keys, matches
	}
}

func createPlayoffRoundsRecursive(playoffPhase models.PlayoffPhase,
	teams []models.TournamentTeamId,
	teamsQuantity int,
	rounds []models.PlayoffRound,
	keys []models.PlayoffRoundKey,
	matches []models.Match) ([]models.PlayoffRound, []models.PlayoffRoundKey, []models.Match) {
	if teamsQuantity <= 1 {
		return rounds, keys, matches
	}

	round := models.PlayoffRound{
		Id:             primitive.NewObjectID(),
		Round:          models.GetRoundFromTeamsCount(teamsQuantity),
		TeamsQuantity:  teamsQuantity,
		PlayoffPhaseId: playoffPhase.Id.Hex(),
	}

	roundKeys := make([]models.PlayoffRoundKey, teamsQuantity/2)
	for i := 0; i < teamsQuantity/2; i++ {
		keyNumber := i + 1
		key := models.PlayoffRoundKey{
			Id:             primitive.NewObjectID(),
			KeyNumber:      strconv.Itoa(keyNumber),
			PlayoffRoundId: round.Id.Hex(),
		}
		roundKeys[i] = key
	}

	if playoffPhase.Config.ClassifiedNumber == 0 { // Solo se crean los partidos si no hay clasificados. Esto quiere decir que es una segunda fase
		if playoffPhase.Config.RandomOrder {
			source := rand.NewSource(time.Now().UnixNano())
			random := rand.New(source)
			random.Shuffle(teamsQuantity, func(i, j int) { teams[i], teams[j] = teams[j], teams[i] })
		}

		if len(rounds) == 0 {
			for i := 0; i < len(roundKeys); i++ {
				roundKeys[i].Teams[0] = teams[i]
				roundKeys[i].Teams[1] = teams[teamsQuantity-1-i]
				roundKeys[i].TeamsRanking[0].TeamId = teams[i]
				roundKeys[i].TeamsRanking[1].TeamId = teams[teamsQuantity-1-i]
			}
			firstRoundMatches := createFirstRoundMatches(playoffPhase, roundKeys)
			matches = append(matches, firstRoundMatches...)
		}
	}

	keys = append(keys, roundKeys...)

	rounds = append(rounds, round)

	halfTeamsCount := teamsQuantity / 2

	if playoffPhase.Config.ClassifiedNumber == 0 {
		return createPlayoffRoundsRecursive(playoffPhase, teams[:halfTeamsCount], halfTeamsCount, rounds, keys, matches)
	} else {
		return createPlayoffRoundsRecursive(playoffPhase, teams, halfTeamsCount, rounds, keys, matches)
	}
}

func createFirstRoundMatches(playoffPhase models.PlayoffPhase, roundKeys []models.PlayoffRoundKey) []models.Match {
	matches := []models.Match{}

	for i := 0; i < len(roundKeys); i++ {
		teamA := roundKeys[i].Teams[0]
		teamB := roundKeys[i].Teams[1]

		match := models.GeneratePlayoffMatch(roundKeys[i].Id.Hex(), teamA, teamB)
		matches = append(matches, match)

		if playoffPhase.Config.HomeAndAway {
			matchReturn := models.GeneratePlayoffMatch(roundKeys[i].Id.Hex(), teamB, teamA)
			matches = append(matches, matchReturn)
		}
	}

	return matches
}
