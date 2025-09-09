package models

import (
	"time"
)

type Iata struct {
	Code        string    `bson:"code" json:"code"`
	Name        string    `bson:"name" json:"name"`
	Latitude    float64   `bson:"latitude" json:"latitude"`
	Longitude   float64   `bson:"longitude" json:"longitude"`
	CountryCode string    `bson:"countryCode" json:"countryCode"`
	ExpiresAt   time.Time `bson:"expiresAt" json:"-"`
}

type IataList struct {
	Data []Iata `json:"data"`
}
