package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarsDto struct {
	CarType string          `bson:"carType" json:"carType"`
	Image   types.FileField `bson:"image" json:"image"`
}

type Cars struct {
	CarsDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type CarsPagination struct {
	Cars []Cars `bson:"cars" json:"cars"`

	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
