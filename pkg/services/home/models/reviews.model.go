package models

import (
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserImage represents the user's profile image
type UserImage struct {
	Path     string `bson:"path" json:"path"`
	Filename string `bson:"filename" json:"filename"`
	Size     int64  `bson:"size" json:"size"`
	Mimetype string `bson:"mimetype" json:"mimetype"`
}

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

// Customer represents the customer information
type Customer struct {
	Name    string     `bson:"name" json:"name"`
	Email   string     `bson:"email,omitempty" json:"email,omitempty"`
	Phone   string     `bson:"phone,omitempty" json:"phone,omitempty"`
	UserImg *UserImage `bson:"userImg,omitempty" json:"userImg,omitempty"`
	UserId  string     `bson:"userId,omitempty" json:"userId,omitempty"`
}

// TravelProgram represents the travel program being reviewed
type TravelProgram struct {
	ProgramId    primitive.ObjectID `bson:"programId,omitempty" json:"programId,omitempty"`
	ProgramTitle string             `bson:"programTitle" json:"programTitle"`
	Destination  string             `bson:"destination,omitempty" json:"destination,omitempty"`
	Countries    []string           `bson:"countries,omitempty" json:"countries,omitempty"`
}

// ReviewDto represents the data transfer object for creating/updating reviews
type ReviewDto struct {
	// Customer type selection (customer or agent)
	Customer string `bson:"customer" json:"customer" validate:"required,oneof=customer agent"`

	// Username from frontend based on auth
	Username string `bson:"username" json:"username" validate:"required"`

	// Program selection (returns only programId)
	ProgramId primitive.ObjectID `bson:"programId" json:"programId" validate:"required"`

	// Independent destination and countries
	Destination string   `bson:"destination" json:"destination"`
	Countries   []string `bson:"countries" json:"countries"`

	// Review content in exact order
	Description        string        `bson:"description" json:"description" validate:"required"`
	AdviceForTravelers string        `bson:"adviceForTravelers" json:"adviceForTravelers"`
	Value              float64       `bson:"value" json:"value" validate:"required,min=1,max=5"`
	Images             []ReviewImage `bson:"images" json:"images"`
	Status             string        `bson:"status" json:"status"` // pending, approved, rejected
	Date               time.Time     `bson:"date" json:"date"`
	Replies            []ReviewReply `bson:"replies" json:"replies"`

	// Legacy fields for backward compatibility (kept for existing data)
	UserImg *UserImage `bson:"userImg,omitempty" json:"userImg,omitempty"`
	UserId  string     `bson:"userId,omitempty" json:"userId,omitempty"`
	Text    string     `bson:"text,omitempty" json:"text,omitempty"`
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

// ReviewWithStats represents a review with additional statistics
type ReviewWithStats struct {
	Review
	ReplyCount int `json:"replyCount"`
	ImageCount int `json:"imageCount"`
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

// ProgramReviewSummary represents review summary for a specific program
type ProgramReviewSummary struct {
	ProgramId       primitive.ObjectID `json:"programId"`
	TotalReviews    int64              `json:"totalReviews"`
	AverageRating   float64            `json:"averageRating"`
	RatingBreakdown map[int]int64      `json:"ratingBreakdown"`
}
