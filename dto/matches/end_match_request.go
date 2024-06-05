package dto

type EndMatchRequest struct {
	AuthorizationCode string `bson:"authorization_code" json:"authorization_code"`
	Comments          string `bson:"comments" json:"comments"`
}
