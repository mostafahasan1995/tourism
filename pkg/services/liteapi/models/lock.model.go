package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lock struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Country   string             `json:"country" bson:"country"`
	Language  string             `json:"language" bson:"language"`
	IsLocked  bool               `json:"isLocked" bson:"isLocked"`
	ExpiresAt time.Time          `json:"expiresAt" bson:"expiresAt"`
}
