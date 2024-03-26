package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchCoach struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchId       string             `bson:"match_id" json:"match_id"`
	TeamId        TournamentTeamId   `bson:"team_id" json:"team_id"`
	CoachId       string             `bson:"coach_id" json:"coach_id"`
	Sanctions     `bson:"sanctions" json:"sanctions"`
	AssociationId string `bson:"association_id" json:"association_id"`
	Status_Data   `bson:"status_data" json:"status_data"`
}

func (matchCoach *MatchCoach) SetCreatedDate() {
	matchCoach.CreatedDate = time.Now()
}

func (matchCoach *MatchCoach) SetModifiedDate() {
	matchCoach.ModifiedDate = time.Now()
}

func (matchCoach *MatchCoach) SetAssociationId(associationId string) {
	matchCoach.AssociationId = associationId
}

func (matchCoach *MatchCoach) SetId(id primitive.ObjectID) {
	matchCoach.Id = id
}
