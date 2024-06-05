package dto

import (
	"time"
)

type ProgramMatchRequest struct {
	Date              time.Time `json:"date"`
	Place             string    `json:"place"`
	AuthorizationCode string    `json:"authorization_code"`
}
