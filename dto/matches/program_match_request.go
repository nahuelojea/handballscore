package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgramMatchRequest struct {
	Date         time.Time          `json:"date"`
	Place        string             `json:"place,omitempty"`
	PlaceId      primitive.ObjectID `json:"place_id,omitempty"`
	StreamingUrl string             `json:"streaming_url,omitempty"`
}
