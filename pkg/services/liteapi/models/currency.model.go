package models

import (
	"time"
)

type Currency struct {
	Code      string    `bson:"code" json:"code"`
	Currency  string    `bson:"currency" json:"currency"`
	Countries []string  `bson:"countries" json:"countries"`
	ExpiresAt time.Time `bson:"expiresAt" json:"-"`
}

type CurrencyList struct {
	Data []Currency `json:"data"`
}
