package models

import (
	"time"
)

type Country struct {
	Code      string    `bson:"code" json:"code"`
	Name      string    `bson:"name" json:"name"`
	ExpiresAt time.Time `bson:"expiresAt" json:"-"`
}

type CountryList struct {
	Data []Country `json:"data"`
}
