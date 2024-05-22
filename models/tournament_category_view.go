package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TournamentCategoryView struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        time.Time          `bson:"end_date" json:"end_date"`
	Status         string             `bson:"status" json:"status"`
	Teams          []TournamentTeamId `bson:"teams" json:"teams"`
	ChampionName   string             `bson:"champion_name" json:"champion_name"`
	ChampionAvatar string             `bson:"champion_avatar" json:"champion_avatar"`
	TournamentId   string             `bson:"tournament_id" json:"tournament_id"`
	CategoryId     string             `bson:"category_id" json:"category_id"`
	AssociationId  string             `bson:"association_id" json:"association_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
}

func (tournamentCategoryView *TournamentCategoryView) SetCreatedDate() {
	tournamentCategoryView.CreatedDate = time.Now()
}

func (tournamentCategoryView *TournamentCategoryView) SetModifiedDate() {
	tournamentCategoryView.ModifiedDate = time.Now()
}

func (tournamentCategoryView *TournamentCategoryView) SetAssociationId(associationId string) {
	tournamentCategoryView.AssociationId = associationId
}

func (tournamentCategoryView *TournamentCategoryView) SetId(id primitive.ObjectID) {
	tournamentCategoryView.Id = id
}
