package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrustedPartnerDto struct {
	Title        string             `bson:"title" json:"title" validate:"required"`
	ProgramTitle string             `bson:"programTitle" json:"programTitle"`
	ProgramId    primitive.ObjectID `bson:"programId" json:"programId"`
	Description  string             `bson:"description" json:"description"`
	Image        types.FileField    `bson:"image" json:"image" validate:"required"`
	URL          string             `bson:"url" json:"url"`
	DisplayOrder int                `bson:"displayOrder" json:"displayOrder"`
	IsActive     bool               `bson:"isActive" json:"isActive"`
}

type TrustedPartner struct {
	TrustedPartnerDto `bson:",inline"`

	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type TrustedPartnerPagination struct {
	Partners   []TrustedPartner  `bson:"partners" json:"partners"`
	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
