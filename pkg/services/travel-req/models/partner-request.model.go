package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PartnerRequestDto `bson:",inline"`
}
