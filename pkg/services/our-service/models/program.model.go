package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Travel Type values:
// Relaxation Trip (رحلة استجمام)
// Adventure (مغامرة)
// Family Trip (رحلة عائلية)
// Romantic Trip (Honeymoon) (رحلة رومانسية – شهر عسل)
// Cultural Trip (رحلة ثقافية)
// Business Trip (رحلة عمل)
// Shopping Trip (رحلة تسوق)
// Wellness or Medical Tourism (رحلة صحية أو استشفائية)

// ✅ عدد الأفراد (Group Size):
// Solo Traveler (فردي)s
// Couple (زوجان)
// Family (عائلة)
// Small Group (مجموعة صغيرة، عادة 4–8 أشخاص)
// Large Group (مجموعة كبيرة، عادة أكثر من 8 أشخاص)

type Program struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProgramDto `bson:",inline"`
	Trash      bool               `bson:"trash" json:"trash"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type ProgramDto struct {
	Title       string             `bson:"title" json:"title" validate:"required"` // program title
	ServiceType enums.ServiceType  `bson:"serviceType" json:"serviceType"`         // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	TravelReqId primitive.ObjectID `bson:"travelReqId" json:"travelReqId"`
	CustomerId  primitive.ObjectID `bson:"customerId" json:"customerId"`
	Status      string             `bson:"status" json:"status" validate:"required,oneof=pending active unactive"`
	Package     primitive.ObjectID `bson:"package" json:"package" validate:"required"`
	ProgramType string             `bson:"programType" json:"programType" validate:"required,oneof=general custom"` // general - custom
	//
	Source      string          `bson:"source" json:"source"`
	Company     string          `bson:"company" json:"company"`         // todo: maybe we need id here
	Coordinator string          `bson:"coordinator" json:"coordinator"` // todo: maybe we need id here
	Purpose     string          `bson:"purpose" json:"purpose"`
	StartDate   time.Time       `bson:"startDate" json:"startDate" validate:"required"`
	EndDate     time.Time       `bson:"endDate" json:"endDate" validate:"required"`
	GroupSize   enums.GroupSize `bson:"groupSize" json:"groupSize" validate:"required,oneof=solo couple family small large"` //see group size values above
	//
	GeneralType *GeneralProgram `bson:"generalType,omitempty" json:"generalType,omitempty" validate:"required_if=ProgramType general"`
	CustomType  *CustomProgram  `bson:"customType,omitempty" json:"customType,omitempty" validate:"required_if=ProgramType custom"`
}

func (p *ProgramDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, p)
}

type ProgramPagination struct {
	Programs   []Program        `bson:"programs" json:"programs"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
