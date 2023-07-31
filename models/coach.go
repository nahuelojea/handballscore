package models

type Coach struct {
	Person
	TeamId string `bson:"team_id" json:"team_id,omitempty"`
}
