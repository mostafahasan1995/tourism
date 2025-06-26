package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewReply struct {
	Text string    `bson:"text" json:"text"`
	Date time.Time `bson:"date" json:"date"`
}

type ReviewDto struct {
	Type               string                 `bson:"type" json:"type" validate:"required"`
	Ref                primitive.ObjectID     `bson:"ref" json:"ref"`
	UserId             string                 `bson:"userId" json:"userId"`
	FirstName          string                 `bson:"firstName" json:"firstName"`
	LastName           string                 `bson:"lastName" json:"lastName"`
	Email              string                 `bson:"email" json:"email"`
	Username           string                 `bson:"username" json:"username"`
	UserImg            *types.FileField       `bson:"userImg,omitempty" json:"userImg,omitempty"`
	Destination        string                 `bson:"destination,omitempty" json:"destination,omitempty"`
	Countries          []string               `bson:"countries,omitempty" json:"countries,omitempty"`
	Description        interface{}            `bson:"description" json:"description" validate:"required"`
	AdviceForTravelers string                 `bson:"adviceForTravelers,omitempty" json:"adviceForTravelers,omitempty"`
	Value              float64                `bson:"value" json:"value" validate:"required,min=1,max=5"`
	Images             []types.FileField      `bson:"images" json:"images"`
	Status             string                 `bson:"status" json:"status"`
	Date               time.Time              `bson:"date" json:"date"`
	Replies            []ReviewReply          `bson:"replies" json:"replies"`
	Customer           string                 `bson:"customer,omitempty" json:"customer,omitempty"`
	Metadata           map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Text               string                 `bson:"text,omitempty" json:"text,omitempty"`
	//ProfileImage       *types.FileField       `bson:"profileImage" json:"profileImage"`
}

func (r *ReviewDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, r)
}

type Review struct {
	ReviewDto `bson:",inline"`
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type ReviewPagination struct {
	Reviews    []Review         `bson:"reviews" json:"reviews"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

type ReviewStats struct {
	TotalReviews      int64         `json:"totalReviews"`
	AverageRating     float64       `json:"averageRating"`
	RatingBreakdown   map[int]int64 `json:"ratingBreakdown"`
	TopDestinations   []string      `json:"topDestinations"`
	TopCountries      []string      `json:"topCountries"`
	ReviewsWithImages int64         `json:"reviewsWithImages"`
}

type EntityReviewSummary struct {
	Type            string             `json:"type"`
	RefId           primitive.ObjectID `json:"refId"`
	TotalReviews    int64              `json:"totalReviews"`
	AverageRating   float64            `json:"averageRating"`
	RatingBreakdown map[int]int64      `json:"ratingBreakdown"`
}

type ReviewStatusDto struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"`
}
