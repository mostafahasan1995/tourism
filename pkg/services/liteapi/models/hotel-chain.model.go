package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelChain struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Id        int                `bson:"id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
