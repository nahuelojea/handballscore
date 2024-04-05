package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchHeaderView struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date           time.Time          `bson:"date" json:"date"`
	TeamHomeName   string             `bson:"team_home_name" json:"team_home_name"`
	TeamHomeAvatar string             `bson:"team_home_avatar" json:"team_home_avatar"`
	TeamAwayName   string             `bson:"team_away_name" json:"team_away_name"`
	TeamAwayAvatar string             `bson:"team_away_avatar" json:"team_away_avatar"`
	Place          string             `bson:"place" json:"place"`
	Status         string             `bson:"status" json:"status"`
	GoalsHome      int                `bson:"goals_home" json:"goals_home"`
	GoalsAway      int                `bson:"goals_away" json:"goals_away"`
	AssociationId  string             `bson:"association_id" json:"association_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
}

func (matchView *MatchHeaderView) SetCreatedDate() {
	matchView.CreatedDate = time.Now()
}

func (matchView *MatchHeaderView) SetModifiedDate() {
	matchView.ModifiedDate = time.Now()
}

func (matchView *MatchHeaderView) SetAssociationId(associationId string) {
	matchView.AssociationId = associationId
}

func (matchView *MatchHeaderView) SetId(id primitive.ObjectID) {
	matchView.Id = id
}
