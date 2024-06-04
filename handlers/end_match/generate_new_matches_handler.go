package end_match

import (
	"errors"
	"strconv"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/matches_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_phases_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_round_keys_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_rounds_repository"
)

type GenerateNewMatchesHandler struct {
	BaseEndMatchHandler
}

func (c *GenerateNewMatchesHandler) HandleEndMatch(endMatch *models.EndMatch) {
	var err error
	status := "Step ignored"

	switch endMatch.CurrentPhase {
	case models.League_Current_Phase:
		err = handleNewMatchesLeaguePhase(endMatch, &status)
	case models.Playoff_Current_Phase:
		err = handleNewMatchesPlayoffPhase(endMatch, &status)
	}

	endMatch.GenerateNewPhase = models.StepStatus{
		IsDone: err == nil,
		Status: func() string {
			if err != nil {
				return err.Error()
			} else {
				return status
			}
		}(),
	}

	if nextHandler := c.GetNext(); nextHandler != nil {
		nextHandler.HandleEndMatch(endMatch)
	}
}

func handleNewMatchesLeaguePhase(endMatch *models.EndMatch, status *string) error {
	leaguePhase := &endMatch.CurrentLeaguePhase.LeaguePhase
	if !leaguePhase.Finished || leaguePhase.Config.ClassifiedNumber == 0 {
		return nil
	}

	leaguePhase.SortTeamsRanking()
	classifiedTeams := leaguePhase.TeamsRanking[:leaguePhase.Config.ClassifiedNumber]

	playoffPhases, _, _, err := getPlayoffPhases(endMatch.CurrentTournamentCategory)
	if err != nil {
		return err
	}

	playoffRounds, _, _, err := getPlayoffRounds(playoffPhases[0].Id.Hex(), leaguePhase.Config.ClassifiedNumber)
	if err != nil {
		return err
	}

	playoffRoundKeys, _, _, err := getPlayoffRoundKeys(playoffRounds[0].Id.Hex())
	if err != nil {
		return err
	}

	assignTeamsToPlayoffKeys(classifiedTeams, playoffRoundKeys)

	if err := updatePlayoffRoundKeys(playoffRoundKeys); err != nil {
		return err
	}

	matches := models.CreateRoundMatches(playoffPhases[0], playoffRoundKeys)
	if _, _, err := matches_repository.CreateMatches(endMatch.CurrentTournamentCategory.AssociationId, matches); err != nil {
		return errors.New("Error to create matches: " + err.Error())
	}

	*status = "New matches generated to end league phase"
	return nil
}

func handleNewMatchesPlayoffPhase(endMatch *models.EndMatch, status *string) error {
	playoffRoundKey := &endMatch.CurrentPlayoffPhase.PlayoffRoundKey

	if !playoffRoundKey.Finished || len(playoffRoundKey.NextRoundKeyId) == 0 {
		return nil
	}

	nextPlayoffRoundKey, _, err := getNextPlayoffRoundKey(playoffRoundKey.NextRoundKeyId)
	if err != nil {
		return err
	}

	playoffRoundKey.SortTeamsRanking()
	classifiedTeam := playoffRoundKey.TeamsRanking[0].TeamId
	firstTeamInKey := isEvenKey(playoffRoundKey.KeyNumber)

	matches, err := matches_repository.GetPendingMatchesByPlayoffRoundKeyId(nextPlayoffRoundKey.Id.Hex())
	if err != nil {
		return errors.New("Error to get pending matches by playoff round key id: " + err.Error())
	}

	if len(matches) == 0 {
		err = createNewMatches(endMatch, nextPlayoffRoundKey, classifiedTeam, firstTeamInKey)
	} else {
		err = updateExistingMatches(matches, classifiedTeam, firstTeamInKey)
	}

	if err != nil {
		return err
	}

	*status = "New matches generated to end playoff phase key"
	return nil
}

func getPlayoffPhases(category models.TournamentCategory) ([]models.PlayoffPhase, int64, int, error) {
	options := playoff_phases_repository.GetPlayoffPhasesOptions{
		TournamentCategoryId: category.Id.Hex(),
		AssociationId:        category.AssociationId,
	}
	return playoff_phases_repository.GetPlayoffPhases(options)
}

func getPlayoffRounds(phaseId string, teamsQuantity int) ([]models.PlayoffRound, int64, int, error) {
	options := playoff_rounds_repository.GetPlayoffRoundsOptions{
		PlayoffPhaseId: phaseId,
		TeamsQuantity:  teamsQuantity,
	}
	return playoff_rounds_repository.GetPlayoffRounds(options)
}

func getPlayoffRoundKeys(roundId string) ([]models.PlayoffRoundKey, int64, int, error) {
	options := playoff_round_keys_repository.GetPlayoffRoundKeysOptions{
		PlayoffRoundId: roundId,
	}
	return playoff_round_keys_repository.GetPlayoffRoundKeys(options)
}

func assignTeamsToPlayoffKeys(classifiedTeams []models.TeamScore, playoffRoundKeys []models.PlayoffRoundKey) {
	n := len(classifiedTeams)
	for i := range playoffRoundKeys {
		playoffRoundKeys[i].Teams[0] = classifiedTeams[i].TeamId
		playoffRoundKeys[i].Teams[1] = classifiedTeams[n-1-i].TeamId
		playoffRoundKeys[i].TeamsRanking[0].TeamId = classifiedTeams[i].TeamId
		playoffRoundKeys[i].TeamsRanking[1].TeamId = classifiedTeams[n-1-i].TeamId
	}
}

func updatePlayoffRoundKeys(keys []models.PlayoffRoundKey) error {
	for _, key := range keys {
		if _, err := playoff_round_keys_repository.UpdatePlayoffRoundKey(key, key.Id.Hex()); err != nil {
			return errors.New("Error to update playoff round key: " + err.Error())
		}
	}
	return nil
}

func getNextPlayoffRoundKey(nextRoundKeyId string) (models.PlayoffRoundKey, bool, error) {
	return playoff_round_keys_repository.GetPlayoffRoundKey(nextRoundKeyId)
}

func isEvenKey(keyNumber string) bool {
	num, _ := strconv.Atoi(keyNumber)
	return num%2 == 0
}

func createNewMatches(endMatch *models.EndMatch, nextPlayoffRoundKey models.PlayoffRoundKey, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	match := generateMatch(endMatch, nextPlayoffRoundKey.Id.Hex(), classifiedTeam, firstTeamInKey)
	if _, _, err := matches_repository.CreateMatch(endMatch.Match.AssociationId, match); err != nil {
		return errors.New("Error to create match: " + err.Error())
	}

	if len(nextPlayoffRoundKey.NextRoundKeyId) > 0 {
		if endMatch.CurrentPlayoffPhase.PlayoffPhase.Config.HomeAndAway {
			return createReturnMatch(endMatch, nextPlayoffRoundKey.NextRoundKeyId, classifiedTeam, firstTeamInKey)
		}
	} else {
		if !endMatch.CurrentPlayoffPhase.PlayoffPhase.Config.SingleMatchFinal && endMatch.CurrentPlayoffPhase.PlayoffPhase.Config.HomeAndAway {
			return createReturnMatch(endMatch, nextPlayoffRoundKey.Id.Hex(), classifiedTeam, firstTeamInKey)
		}
	}
	return nil
}

func generateMatch(endMatch *models.EndMatch, roundKeyId string, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) models.Match {
	if firstTeamInKey {
		return models.GeneratePlayoffMatch(endMatch.CurrentTournamentCategory.TournamentId, roundKeyId, classifiedTeam, models.TournamentTeamId{})
	}
	return models.GeneratePlayoffMatch(endMatch.CurrentTournamentCategory.TournamentId, roundKeyId, models.TournamentTeamId{}, classifiedTeam)
}

func createReturnMatch(endMatch *models.EndMatch, nextRoundKeyId string, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	returnMatch := generateMatch(endMatch, nextRoundKeyId, classifiedTeam, !firstTeamInKey)
	if _, _, err := matches_repository.CreateMatch(endMatch.Match.AssociationId, returnMatch); err != nil {
		return errors.New("Error to create return match: " + err.Error())
	}
	return nil
}

func updateExistingMatches(matches []models.Match, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	if len(matches) == 1 {
		return updateSingleMatch(matches[0], classifiedTeam, firstTeamInKey)
	}
	return updateMultipleMatches(matches, classifiedTeam, firstTeamInKey)
}

func updateSingleMatch(match models.Match, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	if firstTeamInKey {
		if _, err := matches_repository.UpdateHomeTeam(match.Id.Hex(), classifiedTeam); err != nil {
			return errors.New("Error to update home team: " + err.Error())
		}
	} else {
		if _, err := matches_repository.UpdateAwayTeam(match.Id.Hex(), classifiedTeam); err != nil {
			return errors.New("Error to update away team: " + err.Error())
		}
	}
	return nil
}

func updateMultipleMatches(matches []models.Match, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	for i, match := range matches {
		var err error
		if (i == 0 && firstTeamInKey) || (i == 1 && !firstTeamInKey) {
			_, err = matches_repository.UpdateHomeTeam(match.Id.Hex(), classifiedTeam)
		} else {
			_, err = matches_repository.UpdateAwayTeam(match.Id.Hex(), classifiedTeam)
		}
		if err != nil {
			return errors.New("Error to update team: " + err.Error())
		}
	}
	return nil
}
