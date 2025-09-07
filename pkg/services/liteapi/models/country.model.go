package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Country struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Code      string             `bson:"code" json:"code"`
	Name      string             `bson:"name" json:"name"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
