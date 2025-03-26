package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SanctionStatus string

const (
	PendingReview SanctionStatus = "pending_review"
	InProgress    SanctionStatus = "in_progress"
	Completed     SanctionStatus = "completed"
)

type PlayerSanction struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IssueDate        time.Time          `bson:"issue_date" json:"issue_date"`
	EndDate          time.Time          `bson:"end_date" json:"end_date"`
	Description      string             `bson:"description" json:"description"`
	MatchId          string             `bson:"match_id" json:"match_id"`
	PlayerId         string             `bson:"player_id" json:"player_id"`
	MatchSuspensions int                `bson:"match_suspensions" json:"match_suspensions"`
	ServedMatches    []string           `bson:"served_matches" json:"served_matches"`
	SanctionStatus   SanctionStatus     `bson:"sanction_status" json:"sanction_status"`
	AssociationId    string             `bson:"association_id" json:"association_id"`
	Status_Data      `bson:"status_data" json:"status_data"`
}

func (playerSanction *PlayerSanction) SetCreatedDate() {
	playerSanction.CreatedDate = time.Now()
}

func (playerSanction *PlayerSanction) SetModifiedDate() {
	playerSanction.ModifiedDate = time.Now()
}

func (playerSanction *PlayerSanction) SetAssociationId(associationId string) {
	playerSanction.AssociationId = associationId
}

func (playerSanction *PlayerSanction) SetId(id primitive.ObjectID) {
	playerSanction.Id = id
}
