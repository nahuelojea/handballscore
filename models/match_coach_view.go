package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchCoachView struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchId       string             `bson:"match_id" json:"match_id"`
	TeamId        TournamentTeamId   `bson:"team" json:"team"`
	CoachId       string             `bson:"coach_id" json:"coach_id"`
	CoachName     string             `bson:"coach_name" json:"coach_name"`
	CoachSurname  string             `bson:"coach_surname" json:"coach_surname"`
	CoachAvatar   string             `bson:"coach_avatar" json:"coach_avatar"`
	Sanctions     `bson:"sanctions" json:"sanctions"`
	AssociationId string `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
}

func (MatchCoachView *MatchCoachView) SetCreatedDate() {
	MatchCoachView.CreatedDate = time.Now()
}

func (MatchCoachView *MatchCoachView) SetModifiedDate() {
	MatchCoachView.ModifiedDate = time.Now()
}

func (MatchCoachView *MatchCoachView) SetAssociationId(associationId string) {
	MatchCoachView.AssociationId = associationId
}

func (MatchCoachView *MatchCoachView) SetId(id primitive.ObjectID) {
	MatchCoachView.Id = id
}
