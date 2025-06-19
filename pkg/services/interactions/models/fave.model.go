package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FaveType string

const (
	FaveTypeProgram    FaveType = "program"
	FaveTypeHotel      FaveType = "hotel"
	FaveTypeDiary      FaveType = "diary"
	FaveTypeExhibition FaveType = "exhibition"
	FaveTypeAgent      FaveType = "agent"
)

type FaveDto struct {
	Type  FaveType           `bson:"type" json:"type" validate:"required,oneof=program hotel diary exhibition agent"` //program - hotel - diary - exhibition
	RefId primitive.ObjectID `bson:"refId" json:"refId" validate:"required"`
	IsFav bool               `bson:"isFav" json:"isFav"`
}

type Fave struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	FaveDto   `bson:",inline"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type FaveItem struct {
	Fave `bson:",inline"`
	Item any `bson:"item" json:"item"`
}

type FavePagination struct {
	Faves      []Fave           `bson:"faves" json:"faves"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
