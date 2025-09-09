package models

import (
	"time"
)

type HotelChain struct {
	Id        int       `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	ExpiresAt time.Time `bson:"expiresAt" json:"-"`
}

type HotelChainList struct {
	Data []HotelChain `json:"data"`
}
