package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidityUnit string

const (
	ValidityUnitDay   ValidityUnit = "day"
	ValidityUnitMonth ValidityUnit = "month"
	ValidityUnitYear  ValidityUnit = "year"
)

type AttemptsUnit string

const (
	AttemptsUnitDay   AttemptsUnit = "day"
	AttemptsUnitMonth AttemptsUnit = "month"
	AttemptsUnitYear  AttemptsUnit = "year"
)

type MysteryBox struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Name           string             `bson:"name" json:"name"`
	Discount       Discount           `bson:"discount" json:"discount"`
	ValidityPeriod Validity           `bson:"validityPeriod" json:"validityPeriod"`
	ExpiryDate     time.Time          `bson:"expiryDate" json:"expiryDate"`
}

type Discount struct {
	Value int    `bson:"value" json:"value"`
	Unit  string `bson:"unit" json:"unit"`
}

type Validity struct {
	Value int          `bson:"value" json:"value"`
	Unit  ValidityUnit `bson:"unit" json:"unit"`
}

type Game struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string             `bson:"name" json:"name"`
	Boxes    []MysteryBox       `bson:"boxes" json:"boxes"`
	Attempts Attempts           `bson:"attempts" json:"attempts"`
}

type Attempts struct {
	Value int          `bson:"value" json:"value"`
	Unit  AttemptsUnit `bson:"unit" json:"unit"`
}
