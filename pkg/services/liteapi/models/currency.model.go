package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Currency struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Code      string             `bson:"code" json:"code"`
	Currency  string             `bson:"currency" json:"currency"`
	Countries []string           `bson:"countries" json:"countries"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
