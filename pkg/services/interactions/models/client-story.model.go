package models

/*
Example JSON structure for ClientStory:
{
    "_id": "507f1f77bcf86cd799439011",
    "title": "Amazing Trip to Paris",
    "country": "France",
    "city": "Paris",
    "description": "Our unforgettable journey through the City of Light",
    "coverImage": [
        {
            "_id": "507f1f77bcf86cd799439012",
            "originalName": "paris-cover.jpg",
            "path": "/uploads/paris-cover.jpg",
            "service": "file-service",
            "expire": "2024-12-31T23:59:59Z",
            "variants": ["thumbnail", "medium", "large"]
        }
    ],
    "videos": [
        {
            "_id": "507f1f77bcf86cd799439013",
            "originalName": "paris-tour.mp4",
            "path": "/uploads/paris-tour.mp4",
            "service": "file-service",
            "expire": "2024-12-31T23:59:59Z",
            "variants": ["preview", "full"]
        }
    ],
    "trash": false,
    "created_at": "2024-03-20T10:00:00Z",
    "created_by": "507f1f77bcf86cd799439014",
    "updated_at": "2024-03-20T10:00:00Z",
    "updated_by": "507f1f77bcf86cd799439014"
}
*/

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
