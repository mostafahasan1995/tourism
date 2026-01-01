package models

import (
	// "larsa-tourism-microservices/pkg/types"

	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarsDto struct {
	CarType transl.Localizable[string] `bson:"carType" json:"carType"`
	Images  []types.FileField          `bson:"images" json:"images"`
}

//car type

type Cars struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CarsDto   `bson:",inline"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type CarsPagination struct {
	Cars       []Cars           `bson:"cars" json:"cars"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
