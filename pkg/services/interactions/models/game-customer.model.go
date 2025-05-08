package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameCustomer struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId   primitive.ObjectID `bson:"userId" json:"userId"`
	Email    string             `bson:"email" json:"email"`
	BoxId    primitive.ObjectID `bson:"box" json:"box"`
	Discount string             `bson:"discount" json:"discount"`
	Validity string             `bson:"validity" json:"validity"`
	Won      bool               `bson:"won" json:"won"`
	Date     time.Time          `bson:"date" json:"date"`
}

type GameCustomerRes struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Email    string             `bson:"email" json:"email"`
	Attempts int                `bson:"attempts" json:"attempts"`
}
