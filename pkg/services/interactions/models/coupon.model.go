package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    string             `bson:"email" json:"email"`
	BoxId    primitive.ObjectID `bson:"box" json:"box"`
	Code     string             `bson:"code" json:"code"` //auto generated
	Discount Discount           `bson:"discount" json:"discount"`
	Validity Validity           `bson:"validity" json:"validity"`
	Won      bool               `bson:"won" json:"won"`
	Date     time.Time          `bson:"date" json:"date"`
}

type TryBoxDto struct {
	BoxId primitive.ObjectID `bson:"boxId" json:"boxId"`
	Email string             `bson:"email" json:"email"`
}
