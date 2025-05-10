package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientStoryDto struct {
	Title       string            `bson:"title" json:"title"`
	Country     string            `bson:"country" json:"country"`
	City        string            `bson:"city" json:"city"`
	Description string            `bson:"description" json:"description"`
	CoverImage  []types.FileField `bson:"coverImage" json:"CoverImage"`
	Videos      []types.FileField `bson:"videos" json:"videos"`
}

type ClientStory struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ClientStoryDto `bson:",inline"`
	Trash          bool               `bson:"trash" json:"trash"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy      primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy      primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}

type ClientStoryWithPagination struct {
	ClientStories []ClientStory    `bson:"clientStories" json:"clientStories"`
	Pagination    types.Pagination `bson:"pagination" json:"pagination"`
}
