package models

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaguePhase struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Teams                []string           `bson:"teams" json:"teams"`
	HomeAndAway          bool               `bson:"home_and_away" json:"home_and_away"`
	ClassifiedNumber     int                `bson:"classified_number" json:"classified_number"`
	TeamsRanking         []TeamScore        `bson:"teams_ranking" json:"teams_ranking"`
	Status_Data          `bson:"status_data" json:"status_data"`
	TournamentCategoryId string `bson:"tournament_category_id" json:"tournament_category_id"`
	AssociationId        string `bson:"association_id" json:"association_id"`
}

type TeamScore struct {
	TeamId        string `bson:"team" json:"team"`
	Points        int    `bson:"points" json:"points"`
	Matches       int    `bson:"matches" json:"matches"`
	Wins          int    `bson:"wins" json:"wins"`
	Draws         int    `bson:"draws" json:"draws"`
	Losses        int    `bson:"losses" json:"losses"`
	GoalsScored   int    `bson:"goals_scored" json:"goals_scored"`
	GoalsConceded int    `bson:"goals_conceded" json:"goals_conceded"`
}

func (leaguePhase *LeaguePhase) SetAssociationId(associationId string) {
	leaguePhase.AssociationId = associationId
}

func (leaguePhase *LeaguePhase) SetCreatedDate() {
	leaguePhase.CreatedDate = time.Now()
}

func (leaguePhase *LeaguePhase) SetModifiedDate() {
	leaguePhase.ModifiedDate = time.Now()
}

func (leaguePhase *LeaguePhase) SetId(id primitive.ObjectID) {
	leaguePhase.Id = id
}

func (leaguePhase *LeaguePhase) InitializeTeamScores() {
	for _, teamId := range leaguePhase.Teams {
		teamScore := TeamScore{
			TeamId:        teamId,
			Points:        0,
			Matches:       0,
			Wins:          0,
			Draws:         0,
			Losses:        0,
			GoalsScored:   0,
			GoalsConceded: 0,
		}
		leaguePhase.TeamsRanking = append(leaguePhase.TeamsRanking, teamScore)
	}
}

/*func (leaguePhase *LeaguePhase) GenerateMatches() []Match {
	var matches []Match

	if leaguePhase.HomeAndAway {
		for i, teamA := range leaguePhase.Teams {
			for j := i + 1; j < len(leaguePhase.Teams); j++ {
				teamB := leaguePhase.Teams[j]

				matches = append(matches, generateLeagueMatch(leaguePhase.Id.Hex(), teamA, teamB))
				matches = append(matches, generateLeagueMatch(leaguePhase.Id.Hex(), teamB, teamA))
			}
		}
	} else {
		totalTeams := len(leaguePhase.Teams)

		for i := 0; i < totalTeams-1; i++ {
			for j := i + 1; j < totalTeams; j++ {
				var local, visiting string

				if (i+j)%2 == 0 {
					local, visiting = leaguePhase.Teams[i], leaguePhase.Teams[j]
				} else {
					local, visiting = leaguePhase.Teams[j], leaguePhase.Teams[i]
				}

				matches = append(matches, generateLeagueMatch(leaguePhase.Id.Hex(), local, visiting))
			}
		}
	}

	return matches
}*/

func (leaguePhase LeaguePhase) GenerateLeaguePhaseWeeks() []LeaguePhaseWeek {
	totalTeams := len(leaguePhase.Teams)
	var weeks int

	if leaguePhase.HomeAndAway {
		weeks = totalTeams*2 - 2
	} else {
		weeks = totalTeams - 1
	}

	var leaguePhaseWeeks []LeaguePhaseWeek

	for i := 1; i <= weeks; i++ {
		leaguePhaseWeek := LeaguePhaseWeek{
			Number:        i,
			LeaguePhaseId: leaguePhase.Id.Hex(),
		}
		leaguePhaseWeeks = append(leaguePhaseWeeks, leaguePhaseWeek)
	}

	return leaguePhaseWeeks
}

func (leaguePhase *LeaguePhase) GenerateMatches(leaguePhaseWeeks []LeaguePhaseWeek) []Match {
	var matches []Match
	var totalTeams = len(leaguePhase.Teams)

	if leaguePhase.HomeAndAway {
		for i, teamA := range leaguePhase.Teams {
			for j := i + 1; j < totalTeams; j++ {
				teamB := leaguePhase.Teams[j]

				matches = append(matches, generateLeagueMatch(leaguePhase.Id.Hex(), teamA, teamB))
				matches = append(matches, generateLeagueMatch(leaguePhase.Id.Hex(), teamB, teamA))
			}
		}
	} else {
		for i := 0; i < totalTeams-1; i++ {
			for j := i + 1; j < totalTeams; j++ {
				var local, visiting string

				if (i+j)%2 == 0 {
					local, visiting = leaguePhase.Teams[i], leaguePhase.Teams[j]
				} else {
					local, visiting = leaguePhase.Teams[j], leaguePhase.Teams[i]
				}

				matches = append(matches, generateLeagueMatch(leaguePhase.Id.Hex(), local, visiting))
			}
		}
	}

	fmt.Printf("Matches quantity: " + strconv.Itoa(len(matches)))

	rand.Seed(time.Now().UnixNano()) // Inicializar la semilla del generador aleatorio

	for i := 0; i < len(leaguePhaseWeeks); i++ {
		totalMatchesByWeek := int(math.Floor(float64(totalTeams) / 2))

		// Verificar si quedan partidos por asignar
		if totalMatchesByWeek <= 0 {
			break
		}

		// Obtener partidos aleatorios
		matchesToAssign := make([]Match, 0)

		// Usar un mapa para realizar un seguimiento de cuántos partidos ha jugado cada equipo
		teamMatchesCount := make(map[string]int)

		for j := 0; j < totalMatchesByWeek; j++ {
			if len(matches) == 0 {
				break // No quedan partidos
			}

			// Elegir un partido aleatorio
			matchIndex := rand.Intn(len(matches))
			match := matches[matchIndex]

			localTeam := match.TeamLocal
			visitingTeam := match.TeamVisiting

			// Verificar si alguno de los equipos ha alcanzado su límite de partidos en esta fecha
			if teamMatchesCount[localTeam] >= 1 || teamMatchesCount[visitingTeam] >= 1 {
				continue
			}

			// Asignar el partido a la jornada
			matchesToAssign = append(matchesToAssign, match)
			matches = append(matches[:matchIndex], matches[matchIndex+1:]...)

			// Actualizar el recuento de partidos para los equipos involucrados
			teamMatchesCount[localTeam]++
			teamMatchesCount[visitingTeam]++
		}

		// Asignar los partidos a la jornada y reducir totalMatchesByWeek
		totalMatchesByWeek -= len(matchesToAssign)
	}

	return matches
}
