package end_match

import (
	"errors"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/league_phases_repository"
	league_phase_weeks_repository "github.com/nahuelojea/handballscore/repositories/league_phases_weeks_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_phases_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_round_keys_repository"
	"github.com/nahuelojea/handballscore/repositories/playoff_rounds_repository"
	tournaments_repository "github.com/nahuelojea/handballscore/repositories/tournaments_category_repository"
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

func EndMatchChainEvents(match *models.Match) error {
	endMatchHandler := configureChainResponsability()

	endMatch, err := buildEndMatchData(match)
	if err != nil {
		return err
	}

	endMatchHandler.HandleEndMatch(&endMatch)

	return nil
}

func buildEndMatchData(match *models.Match) (models.EndMatch, error) {
	var tournamentCategory models.TournamentCategory
	var err error

	if match.TournamentCategoryId != "" {
		tournamentCategory, _, err = tournaments_repository.GetTournamentCategory(match.TournamentCategoryId)
		if err != nil {
			return models.EndMatch{}, errors.New("Error to get tournament category: " + err.Error())
		}
	}

	endMatch := models.EndMatch{
		Match:                     *match,
		CurrentTournamentCategory: tournamentCategory,
	}

	if match.LeaguePhaseWeekId != "" {
		leaguePhaseWeek, _, err := league_phase_weeks_repository.GetLeaguePhaseWeek(match.LeaguePhaseWeekId)
		if err != nil {
			return models.EndMatch{}, errors.New("Error to get league phase week: " + err.Error())
		}
		leaguePhase, _, err := league_phases_repository.GetLeaguePhase(leaguePhaseWeek.LeaguePhaseId)
		if err != nil {
			return models.EndMatch{}, errors.New("Error to get league phase: " + err.Error())
		}

		endMatch.CurrentLeaguePhase.LeaguePhase = leaguePhase
		endMatch.CurrentLeaguePhase.LeaguePhaseWeek = leaguePhaseWeek
		endMatch.CurrentPhase = models.League_Current_Phase
	}

	if match.PlayoffRoundKeyId != "" {
		playoffRoundKey, _, err := playoff_round_keys_repository.GetPlayoffRoundKey(match.PlayoffRoundKeyId)
		if err != nil {
			return models.EndMatch{}, errors.New("Error to get playoff round key: " + err.Error())
		}
		playoffRound, _, err := playoff_rounds_repository.GetPlayoffRound(playoffRoundKey.PlayoffRoundId)
		if err != nil {
			return models.EndMatch{}, errors.New("Error to get playoff round: " + err.Error())
		}
		playoffPhase, _, err := playoff_phases_repository.GetPlayoffPhase(playoffRound.PlayoffPhaseId)
		if err != nil {
			return models.EndMatch{}, errors.New("Error to get playoff phase: " + err.Error())
		}

		endMatch.CurrentPlayoffPhase.PlayoffPhase = playoffPhase
		endMatch.CurrentPlayoffPhase.PlayoffRound = playoffRound
		endMatch.CurrentPlayoffPhase.PlayoffRoundKey = playoffRoundKey
		endMatch.CurrentPhase = models.Playoff_Current_Phase
	}
	return endMatch, nil
}

func configureChainResponsability() *UpdateTeamsScoreHandler {
	updateTeamsScoreHandler := &UpdateTeamsScoreHandler{}
	/*generateNewPhaseHandler := &end_match.GenerateNewPhaseHandler{}
	generateNewPlayoffRoundKeyHandler := &end_match.GenerateNewPlayoffRoundKeyHandler{}
	updateChampionHandler := &end_match.UpdateChampionHandler{}

	updateTeamsScoreHandler.SetNext(generateNewPhaseHandler)
	generateNewPhaseHandler.SetNext(generateNewPlayoffRoundKeyHandler)
	generateNewPlayoffRoundKeyHandler.SetNext(updateChampionHandler)*/

	return updateTeamsScoreHandler
}
