package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InquiryStatus string

const (
	InquiryStatusReceived    InquiryStatus = "received"
	InquiryStatusUnderReview InquiryStatus = "under review"
	InquiryStatusInProgress  InquiryStatus = "in progress"
	InquiryStatusReplied     InquiryStatus = "replied"
	InquiryStatusClosed      InquiryStatus = "closed"
)

type InquiryReply struct {
	Id              primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	VisitorName     string              `bson:"visitorName" json:"visitorName"`
	UserId          *primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	ResponseMessage string              `bson:"responseMessage" json:"responseMessage" validate:"required"`
	Attachments     []types.FileField   `bson:"attachments" json:"attachments"`
	FollowUpNotes   string              `bson:"followUpNotes,omitempty" json:"followUpNotes,omitempty"`
	Status          InquiryStatus       `bson:"status" json:"status" validate:"required,oneof=replied received 'under review' 'in progress' closed"`
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
	CreatedBy       primitive.ObjectID  `bson:"createdBy" json:"createdBy"`
}

type InquiryDto struct {
	HotelId     primitive.ObjectID     `bson:"hotelId" json:"hotelId" validate:"required"`
	VisitorName string                 `bson:"visitorName" json:"visitorName" validate:"required"`
	Email       string                 `bson:"email" json:"email" validate:"required,email"`
	Phone       string                 `bson:"phone,omitempty" json:"phone,omitempty"`
	Message     string                 `bson:"message" json:"message" validate:"required"`
	Subject     string                 `bson:"subject,omitempty" json:"subject,omitempty"`
	Priority    string                 `bson:"priority" json:"priority" validate:"oneof=low medium high urgent"`
	Source      string                 `bson:"source" json:"source"` // "chat", "contact_form", "email", etc.
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

func (i *InquiryDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, i)
}

type Inquiry struct {
	InquiryDto `bson:",inline"`
	Id         primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	Status     InquiryStatus       `bson:"status" json:"status"`
	Replies    []InquiryReply      `bson:"replies" json:"replies"`
	IsRead     bool                `bson:"isRead" json:"isRead"`
	ReadAt     *time.Time          `bson:"readAt,omitempty" json:"readAt,omitempty"`
	ReadBy     *primitive.ObjectID `bson:"readBy,omitempty" json:"readBy,omitempty"`
	Trash      bool                `bson:"trash" json:"trash"`
	CreatedBy  primitive.ObjectID  `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt  time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy  primitive.ObjectID  `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt  time.Time           `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type InquiryPagination struct {
	Inquiries  []Inquiry        `bson:"inquiries" json:"inquiries"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

type AddInquiryReplyDto struct {
	ResponseMessage string            `bson:"responseMessage" json:"responseMessage" validate:"required"`
	Attachments     []types.FileField `bson:"attachments" json:"attachments"`
	FollowUpNotes   string            `bson:"followUpNotes,omitempty" json:"followUpNotes,omitempty"`
	Status          InquiryStatus     `bson:"status" json:"status" validate:"required,oneof=replied received 'under review' 'in progress' closed"`
}

func (a *AddInquiryReplyDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, a)
}

type UpdateInquiryStatusDto struct {
	Status InquiryStatus `bson:"status" json:"status" validate:"required,oneof=replied received 'under review' 'in progress' closed"`
}

func (u *UpdateInquiryStatusDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, u)
}

type InquiryStats struct {
	TotalInquiries      int64            `json:"totalInquiries"`
	StatusBreakdown     map[string]int64 `json:"statusBreakdown"`
	PriorityBreakdown   map[string]int64 `json:"priorityBreakdown"`
	SourceBreakdown     map[string]int64 `json:"sourceBreakdown"`
	AverageResponseTime float64          `json:"averageResponseTime"` // in hours
	UnreadCount         int64            `json:"unreadCount"`
	TodayInquiries      int64            `json:"todayInquiries"`
	WeeklyTrend         []int64          `json:"weeklyTrend"` // Last 7 days
}

type MarkAsReadDto struct {
	IsRead bool `bson:"isRead" json:"isRead"`
}
