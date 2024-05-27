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

type GetPlayoffRoundsOptions struct {
	PlayoffPhaseId string
	AssociationId  string
	Page           int
	PageSize       int
	SortField      string
	SortOrder      int
}

func GetPlayoffRounds(filterOptions GetPlayoffRoundsOptions) ([]models.PlayoffRound, int64, int, error) {
	filters := playoff_rounds_repository.GetPlayoffRoundsOptions{
		PlayoffPhaseId: filterOptions.PlayoffPhaseId,
		AssociationId:  filterOptions.AssociationId,
		Page:           filterOptions.Page,
		PageSize:       filterOptions.PageSize,
		SortField:      filterOptions.SortField,
		SortOrder:      filterOptions.SortOrder,
	}
	return playoff_rounds_repository.GetPlayoffRounds(filters)
}

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
	var rounds []models.PlayoffRound
	var keys []models.PlayoffRoundKey
	var matches []models.Match

	if playoffPhase.Config.ClassifiedNumber == 0 {
		rounds, keys, matches = createPlayoffRoundsRecursive(playoffPhase, playoffPhase.Teams, len(playoffPhase.Teams), nil, nil, nil)
	} else {
		rounds, keys, matches = createPlayoffRoundsRecursive(playoffPhase, playoffPhase.Teams, playoffPhase.Config.ClassifiedNumber, nil, nil, nil)
	}

	linkKeysBetweenRounds(rounds, keys)
	return rounds, keys, matches
}

func linkKeysBetweenRounds(rounds []models.PlayoffRound, keys []models.PlayoffRoundKey) {
	roundMap := make(map[string][]*models.PlayoffRoundKey)
	for i := range keys {
		roundMap[keys[i].PlayoffRoundId] = append(roundMap[keys[i].PlayoffRoundId], &keys[i])
	}

	for i := 0; i < len(rounds)-1; i++ {
		currentRound := rounds[i]
		nextRoundName, err := models.GetNextRound(currentRound.Round)
		if err != nil {
			// If there is no next round, we skip linking
			continue
		}

		var nextRound *models.PlayoffRound
		for j := range rounds {
			if rounds[j].Round == nextRoundName {
				nextRound = &rounds[j]
				break
			}
		}

		// If next round is not found, skip linking
		if nextRound == nil {
			continue
		}

		currentKeys := roundMap[currentRound.Id.Hex()]
		nextKeys := roundMap[nextRound.Id.Hex()]

		for j := range currentKeys {
			keyNumber, err := strconv.Atoi(currentKeys[j].KeyNumber)
			if err != nil {
				// If there's an error converting keyNumber, skip this key
				continue
			}
			nextKeyIndex := (keyNumber - 1) / 2
			currentKeys[j].NextRoundKeyId = nextKeys[nextKeyIndex].Id.Hex()
		}
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

	if playoffPhase.Config.ClassifiedNumber == 0 { // Only create matches if there are no classified teams. This means it's a secondary phase
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

		match := models.GeneratePlayoffMatch(playoffPhase.TournamentCategoryId, roundKeys[i].Id.Hex(), teamA, teamB)
		matches = append(matches, match)

		if playoffPhase.Config.HomeAndAway {
			matchReturn := models.GeneratePlayoffMatch(playoffPhase.TournamentCategoryId, roundKeys[i].Id.Hex(), teamB, teamA)
			matches = append(matches, matchReturn)
		}
	}

	return matches
}
