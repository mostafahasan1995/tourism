package models

import (
	"encoding/json"
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
	UserImg         *types.FileField    `bson:"userImg,omitempty" json:"userImg,omitempty"`
	Date            string              `bson:"date,omitempty" json:"date,omitempty"`
	Time            string              `bson:"time,omitempty" json:"time,omitempty"`
	InquiryDetails  string              `bson:"inquiryDetails,omitempty" json:"inquiryDetails,omitempty"`
	ResponseMessage string              `bson:"responseMessage" json:"responseMessage" validate:"required"`
	Attachments     []types.FileField   `bson:"attachments" json:"attachments"`
	FollowUpNotes   string              `bson:"followUpNotes,omitempty" json:"followUpNotes,omitempty"`
	Status          InquiryStatus       `bson:"status" json:"status" validate:"required,oneof=replied received 'under review' 'in progress' closed"`
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
	CreatedBy       primitive.ObjectID  `bson:"createdBy" json:"createdBy"`
}

// Custom JSON marshalling for InquiryReply to handle zero ObjectIDs
func (ir InquiryReply) MarshalJSON() ([]byte, error) {
	type Alias InquiryReply
	aux := &struct {
		*Alias
		UserId    *primitive.ObjectID `json:"userId,omitempty"`
		CreatedBy *primitive.ObjectID `json:"createdBy,omitempty"`
	}{
		Alias: (*Alias)(&ir),
	}

	// Handle UserId
	if ir.UserId != nil && !ir.UserId.IsZero() {
		aux.UserId = ir.UserId
	}

	// Handle CreatedBy
	if !ir.CreatedBy.IsZero() {
		aux.CreatedBy = &ir.CreatedBy
	}

	return json.Marshal(aux)
}

type InquiryDto struct {
	ExhibitionId primitive.ObjectID     `bson:"exhibitionId" json:"exhibitionId" validate:"required"`
	HotelId      primitive.ObjectID     `bson:"hotelId" json:"hotelId"`
	VisitorName  string                 `bson:"visitorName" json:"visitorName" validate:"required"`
	Email        string                 `bson:"email" json:"email" validate:"required,email"`
	Phone        types.PhoneNumber      `bson:"phone,omitempty" json:"phone,omitempty"`
	Message      string                 `bson:"message" json:"message" validate:"required"`
	Subject      string                 `bson:"subject,omitempty" json:"subject,omitempty"`
	Priority     string                 `bson:"priority" json:"priority" validate:"oneof=low medium high urgent"`
	Source       string                 `bson:"source" json:"source"` // "chat", "contact_form", "email", etc.
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
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

// Custom JSON marshalling for Inquiry to handle zero ObjectIDs
func (i Inquiry) MarshalJSON() ([]byte, error) {
	type Alias Inquiry
	aux := &struct {
		*Alias
		ExhibitionId *primitive.ObjectID `json:"exhibitionId,omitempty"`
		HotelId      *primitive.ObjectID `json:"hotelId,omitempty"`
		ReadBy       *primitive.ObjectID `json:"readBy,omitempty"`
		CreatedBy    *primitive.ObjectID `json:"createdBy,omitempty"`
		UpdatedBy    *primitive.ObjectID `json:"updatedBy,omitempty"`
	}{
		Alias: (*Alias)(&i),
	}

	// Handle ExhibitionId
	if !i.ExhibitionId.IsZero() {
		aux.ExhibitionId = &i.ExhibitionId
	}

	// Handle HotelId
	if !i.HotelId.IsZero() {
		aux.HotelId = &i.HotelId
	}

	// Handle ReadBy
	if i.ReadBy != nil && !i.ReadBy.IsZero() {
		aux.ReadBy = i.ReadBy
	}

	// Handle CreatedBy
	if !i.CreatedBy.IsZero() {
		aux.CreatedBy = &i.CreatedBy
	}

	// Handle UpdatedBy
	if !i.UpdatedBy.IsZero() {
		aux.UpdatedBy = &i.UpdatedBy
	}

	return json.Marshal(aux)
}

type InquiryPagination struct {
	Inquiries  []Inquiry        `bson:"inquiries" json:"inquiries"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

type AddInquiryReplyDto struct {
	Date            string            `bson:"date,omitempty" json:"date,omitempty"`
	Time            string            `bson:"time,omitempty" json:"time,omitempty"`
	InquiryDetails  string            `bson:"inquiryDetails,omitempty" json:"inquiryDetails,omitempty"`
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

type MarkAsReadDto struct {
	IsRead bool `bson:"isRead" json:"isRead"`
}

func (m *MarkAsReadDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, m)
}

type InquiryStats struct {
	TotalInquiries    int64            `json:"totalInquiries"`
	UnreadInquiries   int64            `json:"unreadInquiries"`
	TodayInquiries    int64            `json:"todayInquiries"`
	StatusBreakdown   map[string]int64 `json:"statusBreakdown"`
	PriorityBreakdown map[string]int64 `json:"priorityBreakdown"`
	SourceBreakdown   map[string]int64 `json:"sourceBreakdown"`
	WeeklyTrend       []int64          `json:"weeklyTrend"`
	ResponseTime      float64          `json:"responseTime"` // Average response time in hours
}
