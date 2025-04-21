package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OurCountryDto struct {
	Name        string            `bson:"name" json:"name"`
	Image       types.FileField   `bson:"image" json:"image"`
	Icon       types.FileField   `bson:"icon" json:"icon"`
	Galeres     []types.FileField `bson:"galeres" json:"galeres"`
	Description string            `bson:"description" json:"description"`
}

type OurCountry struct {
	OurCountryDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type OurCountryPagination struct {
	OurCountry []OurCountry `bson:"ourCountry" json:"ourCountry"`

	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
