package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lock struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Language  string             `json:"language" bson:"language"`
	PlaceId   string             `json:"placeId" bson:"placeId"`
	IsLocked  bool               `json:"isLocked" bson:"isLocked"`
	ExpiresAt time.Time          `json:"expiresAt" bson:"expiresAt"`
}
