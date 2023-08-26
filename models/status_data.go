package models

import (
	"time"
)

type Status_Data struct {
	CreatedDate  time.Time `bson:"created_date" json:"created_date"`
	ModifiedDate time.Time `bson:"modified_date" json:"modified_date"`
	Disabled     bool      `bson:"disabled" json:"disabled"`
}
