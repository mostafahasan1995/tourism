package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReviewImage represents images attached to the review
type ReviewImage struct {
	Path     string `bson:"path" json:"path"`
	Filename string `bson:"filename" json:"filename"`
	Size     int64  `bson:"size" json:"size"`
	Mimetype string `bson:"mimetype" json:"mimetype"`
	Alt      string `bson:"alt,omitempty" json:"alt,omitempty"`
}

// ReviewReply represents a reply to a review
type ReviewReply struct {
	Text string    `bson:"text" json:"text"`
	Date time.Time `bson:"date" json:"date"`
}

// ReviewDto represents the data transfer object for creating/updating reviews
type ReviewDto struct {
	// Review type and reference
	Type string             `bson:"type" json:"type" validate:"required"` // "hotel", "program", "destination", etc.
	Ref  primitive.ObjectID `bson:"ref" json:"ref" validate:"required"`   // Reference ID to the entity being reviewed

	// User information (will be set based on auth)
	UserId    string `bson:"userId" json:"userId"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`

	// Legacy fields for backward compatibility
	Username string           `bson:"username" json:"username"`
	UserImg  *types.FileField `bson:"userImg,omitempty" json:"userImg,omitempty"`

	Destination string   `bson:"destination,omitempty" json:"destination,omitempty"`
	Countries   []string `bson:"countries,omitempty" json:"countries,omitempty"`

	// Review content
	Description        string        `bson:"description" json:"description" validate:"required"`
	AdviceForTravelers string        `bson:"adviceForTravelers,omitempty" json:"adviceForTravelers,omitempty"`
	Value              float64       `bson:"value" json:"value" validate:"required,min=1,max=5"`
	Images             []ReviewImage `bson:"images" json:"images"`
	Status             string        `bson:"status" json:"status"` // pending, approved, rejected
	Date               time.Time     `bson:"date" json:"date"`
	Replies            []ReviewReply `bson:"replies" json:"replies"`

	// Additional fields for flexibility
	Customer string                 `bson:"customer,omitempty" json:"customer,omitempty"` // customer or agent
	Metadata map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // For future extensibility

	// Legacy field for backward compatibility
	Text string `bson:"text,omitempty" json:"text,omitempty"`
}

// Review represents the complete review document
type Review struct {
	ReviewDto `bson:",inline"`

	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// ReviewPagination represents paginated review results
type ReviewPagination struct {
	Reviews    []Review          `bson:"reviews" json:"reviews"`
	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}

// ReviewStats represents overall review statistics
type ReviewStats struct {
	TotalReviews      int64         `json:"totalReviews"`
	AverageRating     float64       `json:"averageRating"`
	RatingBreakdown   map[int]int64 `json:"ratingBreakdown"` // e.g., {5: 10, 4: 5, 3: 2, 2: 1, 1: 0}
	TopDestinations   []string      `json:"topDestinations"`
	TopCountries      []string      `json:"topCountries"`
	ReviewsWithImages int64         `json:"reviewsWithImages"`
}

// EntityReviewSummary represents review summary for a specific entity
type EntityReviewSummary struct {
	Type            string             `json:"type"`
	RefId           primitive.ObjectID `json:"refId"`
	TotalReviews    int64              `json:"totalReviews"`
	AverageRating   float64            `json:"averageRating"`
	RatingBreakdown map[int]int64      `json:"ratingBreakdown"`
}

// ReviewStatusDto for updating review status
type ReviewStatusDto struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"`
}
