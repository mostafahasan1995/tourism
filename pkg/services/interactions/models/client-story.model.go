package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientStoryDto struct {
	Title       transl.Localizable[string] `bson:"title" json:"title"`
	Country     string                     `bson:"country" json:"country"`
	City        string                     `bson:"city" json:"city"`
	Description transl.Localizable[string] `bson:"description" json:"description"`
	CoverImage  types.FileField            `bson:"coverImage" json:"CoverImage"`
	Videos      []types.FileField          `bson:"videos" json:"videos"`
}

type ClientStory struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ClientStoryDto `bson:",inline"`
	Trash          bool               `bson:"trash" json:"trash"`
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedBy      primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	UpdatedBy      primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	IsFav          bool               `bson:"isFav" json:"isFav"`
}

type ClientStoryWithPagination struct {
	ClientStories []ClientStory    `bson:"clientStories" json:"clientStories"`
	Pagination    types.Pagination `bson:"pagination" json:"pagination"`
}
