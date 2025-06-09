package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VisitorDto struct {
	HotelId      primitive.ObjectID     `bson:"hotelId" json:"hotelId" validate:"required"`
	FullName     string                 `bson:"fullName" json:"fullName" validate:"required"`
	Nationality  string                 `bson:"nationality" json:"nationality" validate:"required"`
	Email        string                 `bson:"email" json:"email" validate:"required,email"`
	Phone        string                 `bson:"phone" json:"phone" validate:"required"`
	Interests    []string               `bson:"interests" json:"interests"`
	Tags         []string               `bson:"tags" json:"tags"`
	VisitDate    time.Time              `bson:"visitDate" json:"visitDate"`
	Source       string                 `bson:"source" json:"source"` // "website", "direct", "social_media", etc.
	ReferralCode string                 `bson:"referralCode,omitempty" json:"referralCode,omitempty"`
	Notes        string                 `bson:"notes,omitempty" json:"notes,omitempty"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

func (v *VisitorDto) Validate(validator *validator.Validate) error {
	return helpers.GenericValidation(validator, v)
}

type Visitor struct {
	VisitorDto    `bson:",inline"`
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	LastVisitDate time.Time          `bson:"lastVisitDate" json:"lastVisitDate"`
	VisitCount    int                `bson:"visitCount" json:"visitCount"`
	TotalSpent    float64            `bson:"totalSpent" json:"totalSpent"`
	IsActive      bool               `bson:"isActive" json:"isActive"`
	IsVIP         bool               `bson:"isVIP" json:"isVIP"`
	Trash         bool               `bson:"trash" json:"trash"`
	CreatedBy     primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy     primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt     time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type VisitorPagination struct {
	Visitors   []Visitor        `bson:"visitors" json:"visitors"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

type UpdateVisitorTagsDto struct {
	Tags []string `bson:"tags" json:"tags"`
}

func (u *UpdateVisitorTagsDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, u)
}

type UpdateVisitorInterestsDto struct {
	Interests []string `bson:"interests" json:"interests"`
}

func (u *UpdateVisitorInterestsDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, u)
}

type AddVisitorTagDto struct {
	Tag string `bson:"tag" json:"tag" validate:"required"`
}

func (a *AddVisitorTagDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, a)
}

type AddVisitorInterestDto struct {
	Interest string `bson:"interest" json:"interest" validate:"required"`
}

func (a *AddVisitorInterestDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, a)
}

type VisitorStats struct {
	TotalVisitors        int64            `json:"totalVisitors"`
	ActiveVisitors       int64            `json:"activeVisitors"`
	VIPVisitors          int64            `json:"vipVisitors"`
	NationalityBreakdown map[string]int64 `json:"nationalityBreakdown"`
	SourceBreakdown      map[string]int64 `json:"sourceBreakdown"`
	PopularInterests     []string         `json:"popularInterests"`
	PopularTags          []string         `json:"popularTags"`
	VisitTrends          map[string]int64 `json:"visitTrends"` // Date -> count
	AverageSpent         float64          `json:"averageSpent"`
	ReturnVisitorRate    float64          `json:"returnVisitorRate"`
}

type VisitorActivity struct {
	Id           primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
	VisitorId    primitive.ObjectID     `bson:"visitorId" json:"visitorId"`
	HotelId      primitive.ObjectID     `bson:"hotelId" json:"hotelId"`
	ActivityType string                 `bson:"activityType" json:"activityType"` // "visit", "inquiry", "booking", etc.
	Description  string                 `bson:"description" json:"description"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt    time.Time              `bson:"createdAt" json:"createdAt"`
}

type VisitorActivityPagination struct {
	Activities []VisitorActivity `bson:"activities" json:"activities"`
	Pagination types.Pagination  `bson:"pagination" json:"pagination"`
}

type VisitorWithActivity struct {
	Visitor        `bson:",inline"`
	RecentActivity []VisitorActivity `bson:"recentActivity" json:"recentActivity"`
}

type BulkUpdateVisitorsDto struct {
	VisitorIds []primitive.ObjectID `bson:"visitorIds" json:"visitorIds" validate:"required"`
	Tags       *[]string            `bson:"tags,omitempty" json:"tags,omitempty"`
	Interests  *[]string            `bson:"interests,omitempty" json:"interests,omitempty"`
	IsVIP      *bool                `bson:"isVIP,omitempty" json:"isVIP,omitempty"`
	IsActive   *bool                `bson:"isActive,omitempty" json:"isActive,omitempty"`
}

func (b *BulkUpdateVisitorsDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, b)
}
