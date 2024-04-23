package dto

type SuspendRequest struct {
	Comments string `bson:"comments" json:"comments"`
}
