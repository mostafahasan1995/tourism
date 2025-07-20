package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsletterStatus string

const (
	NewsletterStatusActive   NewsletterStatus = "active"
	NewsletterStatusInactive NewsletterStatus = "inactive"
)

type NewsletterDTO struct {
	Email string `json:"email" bson:"email" validate:"required,email"`
}

type Newsletter struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NewsletterDTO `bson:",inline" json:",inline"`
	Status        NewsletterStatus   `bson:"status" json:"status" validate:"required,oneof=active inactive"`
	UpdatedBy     primitive.ObjectID `bson:"updatedBy" json:"updatedBy,omitempty"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	CreatedBy     primitive.ObjectID `bson:"createdBy" json:"createdBy"`
}
