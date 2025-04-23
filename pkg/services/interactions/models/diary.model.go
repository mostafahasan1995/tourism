package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diary struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name          string             `bson:"name" json:"name"`
	Destination   []string           `bson:"destination" json:"destination"`
	Title         string             `bson:"title" json:"title"`
	FaveActivitie []string           `bson:"faveActivitie" json:"faveActivitie"`
	Bio           string             `bson:"bio" json:"bio"`
	Videos        []types.FileField  `bson:"videos" json:"videos"`
	Images        []types.FileField  `bson:"images" json:"images"`
	Tips          string             `bson:"tips" json:"tips"`
	Description   string             `bson:"description" json:"description"`
	Trash         bool               `bson:"trash" json:"trash"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy     primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}
