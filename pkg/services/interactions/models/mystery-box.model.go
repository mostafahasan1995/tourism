package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MysteryBoxDto struct {
	Title    string `bson:"title" json:"title"`
	Boxes    []Box  `bson:"boxes" json:"boxes"`
	NumOfBox int    `bson:"numOfBox" json:"numOfBox"`
	Attempts int64  `bson:"attempts" json:"attempts"`
	IsActive bool   `bson:"isActive" json:"isActive"`
}

type MysteryBox struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	MysteryBoxDto `bson:",inline"`
	Trash         bool               `bson:"trash" json:"trash"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy     primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}

type Box struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`
	Prize   string             `bson:"prize" json:"prize"`
	Code    string             `bson:"code" json:"code"`
	Program primitive.ObjectID `bson:"program" json:"program"`
}

type BoxTracking struct {
	UserId       primitive.ObjectID `bson:"userId" json:"userId"`
	MysteryBoxId primitive.ObjectID `bson:"mysteryBoxId" json:"mysteryBoxId"`
	Box          *Box               `bson:"box" json:"box"`
	IsWinner     bool               `bson:"isWinner" json:"isWinner"`
	OpenAt       time.Time          `bson:"open_at" json:"open_at"`
}

type MysteryBoxRes struct {
}
