package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameCustomer struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId primitive.ObjectID `bson:"userId" json:"userId"`
	Email  string             `bson:"email" json:"email"`
	Box    MysteryBox         `bson:"box" json:"box"`
	Date   time.Time          `bson:"date" json:"date"`
}
