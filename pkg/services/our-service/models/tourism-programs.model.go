package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourismProgramDto struct {
	Name        string          `bson:"name" json:"name"`
	Destination string          `bson:"destination" json:"destination"`
	TravelType  string          `bson:"travelType" json:"travelType"`
	Duration    int             `bson:"duration" json:"duration"`
	GroupSize   string          `bson:"groupSize" json:"groupSize"`
	Image       types.FileField `bson:"image" json:"image"`
}

type  TourismProgram struct {
	TourismProgramDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}


type TourismProgramPagination struct {
	TourismProgram        []TourismProgram          `bson:"tourismProgram" json:"tourismProgram"`

	Pagination common.Pagination`bson:"pagination" json:"pagination"`
}

