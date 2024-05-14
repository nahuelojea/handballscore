package models

import (
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayoffRoundKey struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	KeyNumber      string             `bson:"key_number" json:"key_number"`
	Teams          []TournamentTeamId `bson:"tournament_team_id" json:"tournament_team_id"`
	TeamsRanking   []TeamScore        `bson:"teams_ranking" json:"teams_ranking"`
	PlayoffRoundId string             `bson:"playoff_round_id" json:"playoff_round_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
	AssociationId  string `bson:"association_id" json:"association_id"`
}

func (playoffRoundKey *PlayoffRoundKey) SortTeamsRanking() {
	sort.SliceStable(playoffRoundKey.TeamsRanking[:], func(i, j int) bool {
		if playoffRoundKey.TeamsRanking[i].Points != playoffRoundKey.TeamsRanking[j].Points {
			return playoffRoundKey.TeamsRanking[i].Points > playoffRoundKey.TeamsRanking[j].Points
		}
		goalDifferenceA := playoffRoundKey.TeamsRanking[i].GoalsScored - playoffRoundKey.TeamsRanking[i].GoalsConceded
		goalDifferenceB := playoffRoundKey.TeamsRanking[j].GoalsScored - playoffRoundKey.TeamsRanking[j].GoalsConceded
		if goalDifferenceA != goalDifferenceB {
			return goalDifferenceA > goalDifferenceB
		}
		if playoffRoundKey.TeamsRanking[i].GoalsScored != playoffRoundKey.TeamsRanking[j].GoalsScored {
			return playoffRoundKey.TeamsRanking[i].GoalsScored > playoffRoundKey.TeamsRanking[j].GoalsScored
		}
		return playoffRoundKey.TeamsRanking[i].GoalsConceded < playoffRoundKey.TeamsRanking[j].GoalsConceded
	})
}

func (playoffRoundKey *PlayoffRoundKey) SetAssociationId(associationId string) {
	playoffRoundKey.AssociationId = associationId
}

func (playoffRoundKey *PlayoffRoundKey) SetCreatedDate() {
	playoffRoundKey.CreatedDate = time.Now()
}

func (playoffRoundKey *PlayoffRoundKey) SetModifiedDate() {
	playoffRoundKey.ModifiedDate = time.Now()
}

func (playoffRoundKey *PlayoffRoundKey) SetId(id primitive.ObjectID) {
	playoffRoundKey.Id = id
}
