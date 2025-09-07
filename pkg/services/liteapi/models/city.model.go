package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type City struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	City      string             `bson:"city" json:"city"`
	Country   string             `bson:"country" json:"country"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
