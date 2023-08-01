package models

import (
	"time"
)

type Status_Data struct {
	CreatedDate  time.Time `bson:"created_date" json:"created_date,omitempty"`
	ModifiedDate time.Time `bson:"modified_date" json:"modified_date,omitempty"`
	Disabled     bool      `bson:"disabled" json:"disabled,omitempty"`
}
