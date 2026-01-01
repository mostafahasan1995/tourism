package models

import (
	"time"
)

type HotelType struct {
	Id        int       `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	ExpiresAt time.Time `bson:"expiresAt" json:"-"`
}

type HotelTypeList struct {
	Data []HotelType `json:"data"`
}
