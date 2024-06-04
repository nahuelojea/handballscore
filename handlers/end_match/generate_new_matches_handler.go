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
			}
			return status
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

	playoffPhaseOptions := playoff_phases_repository.GetPlayoffPhasesOptions{
		TournamentCategoryId: endMatch.CurrentTournamentCategory.Id.Hex(),
		AssociationId:        endMatch.CurrentTournamentCategory.AssociationId,
	}

	playoffPhases, _, _, err := playoff_phases_repository.GetPlayoffPhases(playoffPhaseOptions)
	if err != nil {
		return errors.New("Error to get playoff phases: " + err.Error())
	}

	playoffRoundsOptions := playoff_rounds_repository.GetPlayoffRoundsOptions{
		PlayoffPhaseId: playoffPhases[0].Id.Hex(),
		TeamsQuantity:  leaguePhase.Config.ClassifiedNumber,
	}

	playoffRounds, _, _, err := playoff_rounds_repository.GetPlayoffRounds(playoffRoundsOptions)
	if err != nil {
		return errors.New("Error to get playoff rounds: " + err.Error())
	}

	playoffRoundKeysOption := playoff_round_keys_repository.GetPlayoffRoundKeysOptions{
		PlayoffRoundId: playoffRounds[0].Id.Hex(),
	}

	playoffRoundKeys, _, _, err := playoff_round_keys_repository.GetPlayoffRoundKeys(playoffRoundKeysOption)
	if err != nil {
		return errors.New("Error to get playoff round keys: " + err.Error())
	}

	n := len(classifiedTeams)
	for i := 0; i < len(playoffRoundKeys); i++ {
		playoffRoundKeys[i].Teams[0] = classifiedTeams[i].TeamId
		playoffRoundKeys[i].Teams[1] = classifiedTeams[n-1-i].TeamId
		playoffRoundKeys[i].TeamsRanking[0].TeamId = classifiedTeams[i].TeamId
		playoffRoundKeys[i].TeamsRanking[1].TeamId = classifiedTeams[n-1-i].TeamId
	}

	for i := range playoffRoundKeys {
		_, err = playoff_round_keys_repository.UpdatePlayoffRoundKey(playoffRoundKeys[i], playoffRoundKeys[i].Id.Hex())
		if err != nil {
			return errors.New("Error to update playoff round key: " + err.Error())
		}
	}

	matches := models.CreateRoundMatches(playoffPhases[0], playoffRoundKeys)

	_, _, err = matches_repository.CreateMatches(endMatch.CurrentTournamentCategory.AssociationId, matches)
	if err != nil {
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

	nextPlayoffRoundKey, _, err := playoff_round_keys_repository.GetPlayoffRoundKey(playoffRoundKey.NextRoundKeyId)
	if err != nil {
		return errors.New("Error to get next playoff round key: " + err.Error())
	}

	playoffRoundKey.SortTeamsRanking()
	classifiedTeam := playoffRoundKey.TeamsRanking[0].TeamId
	keyNumber, _ := strconv.Atoi(playoffRoundKey.KeyNumber)
	firstTeamInKey := keyNumber%2 == 0

	matches, err := matches_repository.GetPendingMatchesByPlayoffRoundKeyId(nextPlayoffRoundKey.Id.Hex())
	if err != nil {
		return errors.New("Error to get pending matches by playoff round key id: " + err.Error())
	}

	if len(matches) == 0 {
		createNewMatches(endMatch, nextPlayoffRoundKey, classifiedTeam, firstTeamInKey)
	} else {
		updateExistingMatches(matches, classifiedTeam, firstTeamInKey)
	}

	*status = "New matches generated to end playoff phase key"

	return nil
}

func createNewMatches(endMatch *models.EndMatch, nextPlayoffRoundKey models.PlayoffRoundKey, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	var match models.Match
	if firstTeamInKey {
		match = models.GeneratePlayoffMatch(endMatch.CurrentTournamentCategory.TournamentId, nextPlayoffRoundKey.Id.Hex(), classifiedTeam, models.TournamentTeamId{})
	} else {
		match = models.GeneratePlayoffMatch(endMatch.CurrentTournamentCategory.TournamentId, nextPlayoffRoundKey.Id.Hex(), models.TournamentTeamId{}, classifiedTeam)
	}

	_, _, err := matches_repository.CreateMatch(endMatch.Match.AssociationId, match)
	if err != nil {
		return errors.New("Error to create match: " + err.Error())
	}

	if len(nextPlayoffRoundKey.NextRoundKeyId) > 0 {
		if endMatch.CurrentPlayoffPhase.PlayoffPhase.Config.HomeAndAway {
			createReturnMatch(endMatch, nextPlayoffRoundKey.NextRoundKeyId, classifiedTeam, firstTeamInKey)
		}
	} else {
		if !endMatch.CurrentPlayoffPhase.PlayoffPhase.Config.SingleMatchFinal && endMatch.CurrentPlayoffPhase.PlayoffPhase.Config.HomeAndAway {
			createReturnMatch(endMatch, nextPlayoffRoundKey.Id.Hex(), classifiedTeam, firstTeamInKey)
		}
	}
	return nil
}

func createReturnMatch(endMatch *models.EndMatch, nextRoundKeyId string, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	var returnMatch models.Match
	if firstTeamInKey {
		returnMatch = models.GeneratePlayoffMatch(endMatch.CurrentTournamentCategory.TournamentId, nextRoundKeyId, models.TournamentTeamId{}, classifiedTeam)
	} else {
		returnMatch = models.GeneratePlayoffMatch(endMatch.CurrentTournamentCategory.TournamentId, nextRoundKeyId, classifiedTeam, models.TournamentTeamId{})
	}

	_, _, err := matches_repository.CreateMatch(endMatch.Match.AssociationId, returnMatch)
	if err != nil {
		return errors.New("Error to create return match: " + err.Error())
	}
	return nil
}

func updateExistingMatches(matches []models.Match, classifiedTeam models.TournamentTeamId, firstTeamInKey bool) error {
	if len(matches) == 1 {
		if firstTeamInKey {
			_, err := matches_repository.UpdateHomeTeam(matches[0].Id.Hex(), classifiedTeam)
			if err != nil {
				return errors.New("Error to update home team: " + err.Error())
			}
		} else {
			_, err := matches_repository.UpdateAwayTeam(matches[0].Id.Hex(), classifiedTeam)
			if err != nil {
				return errors.New("Error to update away team: " + err.Error())
			}
		}
	} else {
		if firstTeamInKey {
			_, err := matches_repository.UpdateHomeTeam(matches[0].Id.Hex(), classifiedTeam)
			if err != nil {
				return errors.New("Error to update home team: " + err.Error())
			}

			_, err = matches_repository.UpdateAwayTeam(matches[1].Id.Hex(), classifiedTeam)
			if err != nil {
				return errors.New("Error to update away team: " + err.Error())
			}
		} else {
			_, err := matches_repository.UpdateAwayTeam(matches[0].Id.Hex(), classifiedTeam)
			if err != nil {
				return errors.New("Error to update away team: " + err.Error())
			}

			_, err = matches_repository.UpdateHomeTeam(matches[1].Id.Hex(), classifiedTeam)
			if err != nil {
				return errors.New("Error to update home team: " + err.Error())
			}
		}
	}
	return nil
}
