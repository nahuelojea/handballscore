package dto

import "time"

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	Association  Association `json:"association"`
	Role         string      `json:"role"`
	RoleId       string      `json:"role_id"`
	TeamId       string      `json:"team_id"`
}

type Association struct {
	Id               string    `bson:"_id,omitempty" json:"id"`
	Name             string    `bson:"name" json:"name"`
	DateOfFoundation time.Time `bson:"date_of_foundation" json:"date_of_foundation"`
	Email            string    `bson:"email" json:"email"`
	Avatar           string    `bson:"avatar" json:"avatar"`
	PhoneNumber      string    `bson:"phone_number" json:"phone_number"`
}
