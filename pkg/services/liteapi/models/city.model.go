package models

import (
	"time"
)

type City struct {
	City      string    `bson:"city" json:"city"`
	Country   string    `bson:"country" json:"-"`
	ExpiresAt time.Time `bson:"expiresAt" json:"-"`
}

type CityList struct {
	Data []City ` json:"data"`
}
