package models

import (
	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PartnerRequestDto struct {
	CompanyName     string        `bson:"companyName" json:"companyName"`
	BusinessType    string        `bson:"businessType" json:"businessType"`
	Website         string        `bson:"website" json:"website"`
	CompanyLocation string        `bson:"companyLocation" json:"companyLocation"`
	ContactDetail   ContactDetail `bson:"contactDetail" json:"contactDetail"`
	Message         string        `bson:"message" json:"message"`
}

type ContactDetail struct {
	FullName    string `bson:"fullName" json:"fullName"`
	Position    string `bson:"position" json:"position"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
	Email       string `bson:"email" json:"email"`
}

type PartnerRequest struct {
	PartnerRequestDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type PartnerRequestPagination struct {
	PartnerRequest []PartnerRequest `bson:"partnerRequest" json:"partnerRequest"`

	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
