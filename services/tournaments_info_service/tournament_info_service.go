package tournaments_info_service

import (
	"errors"
	"strconv"
	"strings"

	TournamentCategoryDTO "github.com/nahuelojea/handballscore/dto/tournament_categories"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/services/league_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_phases_service"
	"github.com/nahuelojea/handballscore/services/playoff_round_keys_service"
	"github.com/nahuelojea/handballscore/services/playoff_rounds_service"
	"github.com/nahuelojea/handballscore/services/teams_service"
	tournaments_service "github.com/nahuelojea/handballscore/services/tournaments_category_service"
)

func GetInfo(id string) (TournamentCategoryDTO.TournamentInfoResponse, error) {
	var tournamentInfoResponse TournamentCategoryDTO.TournamentInfoResponse

	tournamentCategory, _, err := tournaments_service.GetTournamentCategory(id)
	if err != nil {
		return tournamentInfoResponse, errors.New("Error al obtener la categoría del torneo: " + err.Error())
	}

	tournamentInfoResponse.LeaguePhaseInfo, err = getLeaguePhaseInfo(tournamentCategory)
	if err != nil {
		return tournamentInfoResponse, err
	}

	tournamentInfoResponse.PlayoffPhaseInfo, err = getPlayoffPhaseInfo(tournamentCategory)
	if err != nil {
		return tournamentInfoResponse, err
	}

	return tournamentInfoResponse, nil
}

func getLeaguePhaseInfo(tournamentCategory models.TournamentCategory) (TournamentCategoryDTO.LeaguePhaseInfoResponse, error) {
	var leaguePhaseInfo TournamentCategoryDTO.LeaguePhaseInfoResponse

	filterLeaguePhase := league_phases_service.GetLeaguePhasesOptions{
		TournamentCategoryId: tournamentCategory.Id.Hex(),
		AssociationId:        tournamentCategory.AssociationId,
	}

	leaguePhases, _, _, _ := league_phases_service.GetLeaguePhases(filterLeaguePhase)
	if len(leaguePhases) > 0 {
		leaguePhase := leaguePhases[0]

		league_phases_service.ApplyOlympicTiebreaker(&leaguePhase)

		classifiedNumber := leaguePhase.Config.ClassifiedNumber
		teamsRanking := make([]TournamentCategoryDTO.TeamScoreResponse, 0)

		for i, teamScore := range leaguePhase.TeamsRanking {
			teamName := ""
			teamAvatar := ""

			position := i + 1
			classified := position <= classifiedNumber

			team, _, _ := teams_service.GetTeam(teamScore.TeamId.TeamId)
			if len(team.Name) > 0 {
				teamName = team.Name
			}
			if len(team.Avatar) > 0 {
				teamAvatar = team.Avatar
			}

			teamInfo := TournamentCategoryDTO.TeamInfoResponse{
				TeamId:     teamScore.TeamId.TeamId,
				TeamVariant: teamScore.TeamId.Variant,
				TeamName:   strings.TrimSpace(teamName + " " + teamScore.TeamId.Variant),
				TeamAvatar: teamAvatar,
			}

			teamRanking := TournamentCategoryDTO.TeamScoreResponse{
				TeamInfo:        teamInfo,
				Position:        position,
				Classified:      classified,
				Points:          teamScore.Points,
				Matches:         teamScore.Matches,
				Wins:            teamScore.Wins,
				Draws:           teamScore.Draws,
				Losses:          teamScore.Losses,
				GoalsScored:     teamScore.GoalsScored,
				GoalsConceded:   teamScore.GoalsConceded,
				GoalsDifference: teamScore.GoalsScored - teamScore.GoalsConceded,
			}

			teamsRanking = append(teamsRanking, teamRanking)
		}

		leaguePhaseInfo = TournamentCategoryDTO.LeaguePhaseInfoResponse{
			TeamsRanking: teamsRanking,
		}
	}

	return leaguePhaseInfo, nil
}

func getPlayoffPhaseInfo(tournamentCategory models.TournamentCategory) (TournamentCategoryDTO.PlayoffPhaseInfoResponse, error) {
	var playoffPhaseInfo TournamentCategoryDTO.PlayoffPhaseInfoResponse

	filterPlayoffPhase := playoff_phases_service.GetPlayoffPhasesOptions{
		TournamentCategoryId: tournamentCategory.Id.Hex(),
		AssociationId:        tournamentCategory.AssociationId,
	}

	playoffPhases, _, _, _ := playoff_phases_service.GetPlayoffPhases(filterPlayoffPhase)
	if len(playoffPhases) > 0 {
		playoffPhase := playoffPhases[0]

		playoffKeys, err := getPlayoffRoundsInfo(playoffPhase)
		if err != nil {
			return playoffPhaseInfo, err
		}

		playoffPhaseConfig := TournamentCategoryDTO.PlayoffConfigResponse{
			HomeAndAway:      playoffPhase.Config.HomeAndAway,
			SingleMatchFinal: playoffPhase.Config.SingleMatchFinal,
			RandomOrder:      playoffPhase.Config.RandomOrder,
			ClassifiedNumber: playoffPhase.Config.ClassifiedNumber,
		}

		playoffPhaseInfo = TournamentCategoryDTO.PlayoffPhaseInfoResponse{
			PlayoffKeys:   playoffKeys,
			PlayoffConfig: playoffPhaseConfig,
		}
	}

	return playoffPhaseInfo, nil
}

func getPlayoffRoundsInfo(playoffPhase models.PlayoffPhase) ([]TournamentCategoryDTO.PlayoffKeyResponse, error) {
	var playoffRoundsInfo []TournamentCategoryDTO.PlayoffKeyResponse

	filterPlayoffRound := playoff_rounds_service.GetPlayoffRoundsOptions{
		PlayoffPhaseId: playoffPhase.Id.Hex(),
		AssociationId:  playoffPhase.AssociationId,
	}

	playoffRounds, _, _, _ := playoff_rounds_service.GetPlayoffRounds(filterPlayoffRound)

	for _, round := range playoffRounds {
		roundKeysInfo, err := getPlayoffRoundKeysInfo(round)
		if err != nil {
			return nil, errors.New("Error to get playoff round keys: " + err.Error())
		}
		playoffRoundsInfo = append(playoffRoundsInfo, roundKeysInfo...)
	}

	return playoffRoundsInfo, nil
}

func getPlayoffRoundKeysInfo(playoffRound models.PlayoffRound) ([]TournamentCategoryDTO.PlayoffKeyResponse, error) {
	var playoffRoundKeysInfo []TournamentCategoryDTO.PlayoffKeyResponse

	filterPlayoffRoundKey := playoff_round_keys_service.GetPlayoffRoundKeysOptions{
		PlayoffRoundId: playoffRound.Id.Hex(),
		AssociationId:  playoffRound.AssociationId,
	}

	playoffRoundKeys, _, _, _ := playoff_round_keys_service.GetPlayoffRoundKeys(filterPlayoffRoundKey)

	for _, playoffRoundKey := range playoffRoundKeys {
		teamHomeName := ""
		teamHomeAvatar := ""
		teamAwayName := ""
		teamAwayAvatar := ""

		teamHome, _, _ := teams_service.GetTeam(playoffRoundKey.Teams[0].TeamId)
		if len(teamHome.Name) > 0 {
			teamHomeName = teamHome.Name
		}
		if len(teamHome.Avatar) > 0 {
			teamHomeAvatar = teamHome.Avatar
		}

		teamAway, _, _ := teams_service.GetTeam(playoffRoundKey.Teams[1].TeamId)
		if len(teamAway.Name) > 0 {
			teamAwayName = teamAway.Name
		}
		if len(teamAway.Avatar) > 0 {
			teamAwayAvatar = teamAway.Avatar
		}

		name := playoffRound.PlayoffRoundNameTraduction()
		if playoffRound.Round != "final" {
			name += " " + playoffRoundKey.KeyNumber
		}

		isWinnerHome := playoffRoundKey.Winner == playoffRoundKey.Teams[0]
		isWinnerAway := playoffRoundKey.Winner == playoffRoundKey.Teams[1]

		matchStatus := "NO_PARTY"
		var teamsStatus string
		if playoffRoundKey.Winner != (models.TournamentTeamId{}) {
			matchStatus = "DONE"
			teamsStatus = "PLAYED"
		}

		var homeResult string
		var awayResult string

		if len(playoffRoundKey.MatchResults) > 0 {
			var homeResults []string
			var awayResults []string
			for _, matchResult := range playoffRoundKey.MatchResults {
				if matchResult.TeamHome == playoffRoundKey.Teams[0] {
					homeResults = append(homeResults, strconv.Itoa(matchResult.TeamHomeGoals))
					awayResults = append(awayResults, strconv.Itoa(matchResult.TeamAwayGoals))
				} else if matchResult.TeamHome == playoffRoundKey.Teams[1] {
					homeResults = append(homeResults, strconv.Itoa(matchResult.TeamAwayGoals))
					awayResults = append(awayResults, strconv.Itoa(matchResult.TeamHomeGoals))
				}
			}
			homeResult = strings.Join(homeResults, "-")
			awayResult = strings.Join(awayResults, "-")
		}

		playoffKeyTeams := []TournamentCategoryDTO.PlayoffKeyTeamResponse{}

		if playoffRoundKey.Teams[0].TeamId != "" {
			playoffKeyTeams = append(playoffKeyTeams, TournamentCategoryDTO.PlayoffKeyTeamResponse{
				Id: playoffRoundKey.Teams[0].TeamId,
				TeamInfoResponse: TournamentCategoryDTO.TeamInfoResponse{
					TeamName:   teamHomeName + " " + playoffRoundKey.Teams[0].Variant,
					TeamAvatar: teamHomeAvatar,
				},
				Result:   homeResult,
				Status:   teamsStatus,
				IsWinner: isWinnerHome,
			})
		}

		if playoffRoundKey.Teams[1].TeamId != "" {
			playoffKeyTeams = append(playoffKeyTeams, TournamentCategoryDTO.PlayoffKeyTeamResponse{
				Id: playoffRoundKey.Teams[1].TeamId,
				TeamInfoResponse: TournamentCategoryDTO.TeamInfoResponse{
					TeamName:   teamAwayName + " " + playoffRoundKey.Teams[1].Variant,
					TeamAvatar: teamAwayAvatar,
				},
				Result:   awayResult,
				Status:   teamsStatus,
				IsWinner: isWinnerAway,
			})
		}

		playoffRoundKeysInfo = append(playoffRoundKeysInfo, TournamentCategoryDTO.PlayoffKeyResponse{
			Id:               playoffRoundKey.Id.Hex(),
			Name:             name,
			NextPlayoffKeyId: playoffRoundKey.NextRoundKeyId,
			State:            matchStatus,
			PlayoffKeyTeams:  playoffKeyTeams,
		})
	}

	return playoffRoundKeysInfo, nil
}
