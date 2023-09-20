package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coach struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Personal_Data `bson:"personal_data" json:"personal_data"`
	Status_Data   `bson:"status_data" json:"status_data"`
	TeamId        string `bson:"team_id" json:"team_id"`
	AssociationId string `bson:"association_id" json:"association_id"`
}

func (coach *Coach) SetCreatedDate() {
	coach.CreatedDate = time.Now()
}

func (coach *Coach) SetModifiedDate() {
	coach.ModifiedDate = time.Now()
}

func (coach *Coach) SetDisabled(disabled bool) {
	coach.Disabled = disabled
}

func (coach *Coach) SetAssociationId(associationId string) {
	coach.AssociationId = associationId
}

func (coach *Coach) SetAvatarURL(filename string) {
	coach.Avatar = ImagesBaseURL + filename
}
