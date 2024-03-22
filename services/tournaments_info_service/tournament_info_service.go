package tournaments_info_service

import (
	"errors"
	"strings"

	TournamentCategoryDTO "github.com/nahuelojea/handballscore/dto/tournament_categories"
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

func getLeaguePhaseInfo(tournamentCategory TournamentCategoryDTO.TournamentCategory) (TournamentCategoryDTO.LeaguePhaseInfoResponse, error) {
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

func getPlayoffPhaseInfo(tournamentCategory TournamentCategoryDTO.TournamentCategory) (TournamentCategoryDTO.PlayoffPhaseInfoResponse, error) {
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
			PlayoffRounds: playoffRounds,
		}
	}

	return playoffPhaseInfo, nil
}

func getPlayoffRoundsInfo(playoffPhase TournamentCategoryDTO.PlayoffPhase) ([]TournamentCategoryDTO.PlayoffRoundInfoResponse, error) {
	var playoffRoundsInfo []TournamentCategoryDTO.PlayoffRoundInfoResponse

	filterPlayoffRound := playoff_rounds_service.GetPlayoffRoundsOptions{
		PlayoffPhaseId: playoffPhase.Id.Hex(),
		AssociationId:  playoffPhase.AssociationId,
	}

	playoffRounds, _, _, _ := playoff_rounds_service.GetPlayoffRounds(filterPlayoffRound)

	for _, round := range playoffRounds {
		playoffRoundKeysInfo, err := getPlayoffRoundKeysInfo(round)
		if err != nil {
			return nil, err
		}

		playoffRoundInfo := TournamentCategoryDTO.PlayoffRoundInfoResponse{
			Round:               round.PlayoffRoundNameTraduction(),
			PlayoffRoundKeyInfo: playoffRoundKeysInfo,
		}

		playoffRoundsInfo = append(playoffRoundsInfo, playoffRoundInfo)
	}

	return playoffRoundsInfo, nil
}

func getPlayoffRoundKeysInfo(playoffRound TournamentCategoryDTO.PlayoffRound) ([]TournamentCategoryDTO.PlayoffRoundKeyInfoResponse, error) {
	var playoffRoundKeysInfo []TournamentCategoryDTO.PlayoffRoundKeyInfoResponse

	filterPlayoffRoundKey := playoff_round_keys_service.GetPlayoffRoundKeysOptions{
		PlayoffRoundId: playoffRound.Id.Hex(),
		AssociationId:  playoffRound.AssociationId,
	}

	playoffRoundKeys, _, _, _ := playoff_round_keys_service.GetPlayoffRoundKeys(filterPlayoffRoundKey)

	for _, playoffRoundKey := range playoffRoundKeys {
		teamsRanking := make([]TournamentCategoryDTO.TeamScoreResponse, 0)

		playoffRoundKey.SortTeamsRanking()

		for i, teamScore := range playoffRoundKey.TeamsRanking {
			teamName := ""
			teamAvatar := ""

			position := i + 1
			classified := position == 1

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

		playoffRoundKeysInfo = append(playoffRoundKeysInfo, TournamentCategoryDTO.PlayoffRoundKeyInfoResponse{
			TeamsRanking: teamsRanking,
		})
	}

	return playoffRoundKeysInfo, nil
}
