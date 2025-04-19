package models

import (
	// "larsa-tourism-microservices/pkg/types"
	// "larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactUsDto struct {
	FullName        string `bson:"fullName" json:"fullName"`
	EmailAddress    string `bson:"emailAddress" json:"emailAddress"`
	PhoneNumber     string `bson:"phoneNumber" json:"phoneNumber"`
	HowDidYouFindUs string `bson:"howDidYouFindUs" json:"howDidYouFindUs"`
	Message         string `bson:"message" json:"message"`
}

type ContactUs struct {
	ContactUsDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type ContactUsPagination struct {
	ContactUs []ContactUs `bson:"contactUs" json:"contactUs"`

	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
