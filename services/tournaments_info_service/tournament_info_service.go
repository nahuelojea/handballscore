package tournaments_info_service

import (
	"errors"
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

// GetInfo devuelve la información del torneo correspondiente al ID proporcionado.
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

		leaguePhase.SortTeamsRanking()

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

		playoffRounds, err := getPlayoffRoundsInfo(playoffPhase)
		if err != nil {
			return playoffPhaseInfo, err
		}

		playoffPhaseInfo = TournamentCategoryDTO.PlayoffPhaseInfoResponse{
			PlayoffKeys: playoffRounds,
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
		playoffRoundsInfo, err := getPlayoffRoundKeysInfo(round)
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

		team, _, _ := teams_service.GetTeam(playoffRoundKey.Teams[0].TeamId)
		if len(team.Name) > 0 {
			teamHomeName = team.Name
		}
		if len(team.Avatar) > 0 {
			teamHomeName = team.Avatar
		}

		team, _, _ = teams_service.GetTeam(playoffRoundKey.Teams[1].TeamId)
		if len(team.Name) > 0 {
			teamAwayName = team.Name
		}
		if len(team.Avatar) > 0 {
			teamAwayName = team.Avatar
		}

		playoffRoundKeysInfo = append(playoffRoundKeysInfo, TournamentCategoryDTO.PlayoffKeyResponse{
			Id:               playoffRoundKey.Id.Hex(),
			Name:             playoffRound.Round + " - Key " + playoffRoundKey.KeyNumber,
			NextPlayoffKeyId: playoffRoundKey.NextPlayoffKeyId.Hex(),
			State:            playoffRoundKey.State,
			PlayoffKeyTeams: []TournamentCategoryDTO.PlayoffKeyTeamResponse{
				{
					Id: playoffRoundKey.Teams[0].TeamId.Hex(),
					TeamInfoResponse: TournamentCategoryDTO.TeamInfoResponse{
						TeamName:   teamHomeName,
						TeamAvatar: teamHomeAvatar,
					},
					Result: playoffRoundKey.Teams[0].Result,
					Status: playoffRoundKey.Teams[0].Status,
				},
				{
					Id: playoffRoundKey.Teams[1].TeamId.Hex(),
					TeamInfoResponse: TournamentCategoryDTO.TeamInfoResponse{
						TeamName:   teamAwayName,
						TeamAvatar: teamAwayAvatar,
					},
					Result: playoffRoundKey.Teams[1].Result,
					Status: playoffRoundKey.Teams[1].Status,
				},
			},
		})
	}

	return playoffRoundKeysInfo, nil
}
