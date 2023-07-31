package models

type Player struct {
	Person
	Gender          string `bson:"gender" json:"gender,omitempty"`
	AffiliateNumber string `bson:"affiliate_number" json:"affiliate_number,omitempty"`
	TeamId          string `bson:"team_id" json:"team_id,omitempty"`
}
