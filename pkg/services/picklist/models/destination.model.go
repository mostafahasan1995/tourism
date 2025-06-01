package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DestinationDto struct {
	Name        string            `bson:"name" json:"name"`
	Images      []types.FileField `bson:"images" json:"images"`
	Icon        types.FileField   `bson:"icon" json:"icon"`
	Description string            `bson:"description" json:"description"`
}

type Destination struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DestinationDto `bson:",inline"`
	Trash          bool               `bson:"trash" json:"trash"`
	CreatedBy      primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy      primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type DestinationPagination struct {
	Destinations []Destination    `bson:"destinations" json:"destinations"`
	Pagination   types.Pagination `bson:"pagination" json:"pagination"`
}
