package dto

import (
	"time"
)

type ProgramMatchRequest struct {
	Date         time.Time `json:"date"`
	Place        string    `json:"place,omitempty"`
	PlaceId      string    `json:"place_id,omitempty"`
	StreamingUrl string    `json:"streaming_url,omitempty"`
}
