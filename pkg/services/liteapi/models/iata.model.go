package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Iata struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Latitude    float64            `bson:"latitude" json:"latitude"`
	Longitude   float64            `bson:"longitude" json:"longitude"`
	CountryCode string             `bson:"countryCode" json:"countryCode"`
	ExpiresAt   time.Time          `bson:"expiresAt" json:"expiresAt"`
}
