package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Program struct {
	BasicInfo      `bson:",inline"`
	GeneralProgram *GeneralProgram    `bson:"generalProgram,omitempty" json:"generalProgram,omitempty"`
	CustomProgram  *CustomProgram     `bson:"customProgram,omitempty" json:"customProgram,omitempty"`
	Trash          bool               `bson:"trash" json:"trash"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy      primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy      primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}
