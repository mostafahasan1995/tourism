package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerDto struct {
	CustomerName        string            `bson:"customerName" json:"customerName"`
	Nationality         string            `bson:"nationality" json:"nationality"`
	Company             string            `bson:"company" json:"company"`
	About               string            `bson:"about" json:"about"`
	Photo               []types.FileField `bson:"photo" json:"photo"`
	ClientMobile        string            `bson:"clientMobile" json:"clientMobile"`
	ClientWhatsapp      string            `bson:"clientWhatsapp" json:"clientWhatsapp"`
	ClientEmail         string            `bson:"clientEmail" json:"clientEmail"`
	CoordinatorMobile   string            `bson:"coordinatorMobile" json:"coordinatorMobile"`
	CoordinatorWhatsapp string            `bson:"coordinatorWhatsapp" json:"coordinatorWhatsapp"`
	CoordinatorEmail    string            `bson:"coordinatorEmail" json:"coordinatorEmail"`
	SocialMedia         []SocialMedia     `bson:"socialMedia" json:"socialMedia"`
	Email               string            `bson:"email" json:"email"`
}

type SocialMedia struct {
	Platform string `bson:"platform" json:"platform"` //e.g. facebook, instagram, twitter, linkedin, etc.
	Link     string `bson:"link" json:"link"`
}

type Customer struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CustomerId  string             `bson:"customerId,omitempty" json:"customerId,omitempty"`
	CustomerDto `bson:",inline"`
	Trash       bool               `bson:"trash" json:"trash"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy   primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}
