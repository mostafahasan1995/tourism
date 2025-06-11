package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhoneNumber struct {
	Prefix  string `bson:"prefix" json:"prefix" validate:"required"`   // Country code
	Content string `bson:"content" json:"content" validate:"required"` // Actual phone number
}

type PartnerRequestDto struct {
	// Business Information
	CompanyName     string `bson:"companyName" json:"companyName" validate:"required"`
	BusinessType    string `bson:"businessType" json:"businessType" validate:"required"`
	Website         string `bson:"website" json:"website"`
	CompanyLocation string `bson:"companyLocation" json:"companyLocation" validate:"required"`

	// Contact Details
	FullName    string      `bson:"fullName" json:"fullName" validate:"required"`
	Position    string      `bson:"position" json:"position" validate:"required"`
	PhoneNumber PhoneNumber `bson:"phoneNumber" json:"phoneNumber" validate:"required"`
	Email       string      `bson:"email" json:"email" validate:"required,email"`

	// Services Offered
	ServicesOffered []string `bson:"servicesOffered" json:"servicesOffered" validate:"required"`

	// Additional Fields
	CompanyLogo types.FileField `bson:"companyLogo" json:"companyLogo"`
	Notes       string          `bson:"notes" json:"notes"`
	Status      string          `bson:"status" json:"status"` // "pending", "approved", "rejected"
}

type PartnerRequest struct {
	PartnerRequestDto `bson:",inline"`

	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
