package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	membermodels "larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelRequestDto struct {
	//basic information
	//BasicInfo BasicInfo `bson:"basicInfo" json:"basicInfo"`

	ClientName   string `bson:"clientName" json:"clientName" validate:"required"`
	ClientPhone  string `bson:"clientPhone" json:"clientPhone" `
	ClientEmail  string `bson:"clientEmail" json:"clientEmail" validate:"required"`
	Nationality  string `bson:"nationality" json:"nationality" `
	TripDuration int    `bson:"tripDuration" json:"tripDuration" validate:"required"`
	//
	ServiceType enums.ServiceType `bson:"serviceType" json:"serviceType" validate:"required,oneof=delegation custom-plan business-man vip-car flight-request partner-request hotel-booking relaxation adventure family romantic cultural business shopping wellness"` // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	//request
	Delegation          *Delegation          `bson:"delegation,omitempty" json:"delegation,omitempty" validate:"required_if=ServiceType delegation"`
	BusinessMan         *BusinessMan         `bson:"businessMan,omitempty" json:"businessMan,omitempty" validate:"required_if=ServiceType business-man"`
	CustomPlan          *CustomPlan          `bson:"customPlan,omitempty" json:"customPlan,omitempty" validate:"required_if=ServiceType custom-plan"`
	VipCar              *VipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty" validate:"required_if=ServiceType vip-car"`
	FlightTicketRequest *FlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty" validate:"required_if=ServiceType flight-request"`
	PartnerRequest      *PartnerRequest      `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty" validate:"required_if=ServiceType partner-request"`
	//delegation - business-man - custom-plan info
	Destination     []Destination      `bson:"destination,omitempty" json:"destination,omitempty" `
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator" validate:"required"`
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod" validate:"required"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}

// type BasicInfo struct {
// 	ProgramTitle   string     `bson:"programTitle" json:"programTitle"`
// 	ServiceType    string     `bson:"serviceType" json:"serviceType"`
// 	Purpose        string     `bson:"purpose" json:"purpose"`
// 	DelegationType string     `bson:"delegationType" json:"delegationType"`
// 	Source         string     `bson:"source" json:"source"`
// 	Customer       string     `bson:"customer" json:"customer"`
// 	Company        string     `bson:"company" json:"company"`
// 	Coordinator    string     `bson:"coordinator" json:"coordinator"`
// 	StartDate      *time.Time `bson:"startDate" json:"startDate"`
// 	EndDate        *time.Time `bson:"endDate" json:"endDate"`
// 	GroupSize      string     `bson:"groupSize" json:"groupSize"`
// }

func (t *TravelRequestDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, t)
}

type TravelRequest struct {
	Id               primitive.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	ReqId            string                `bson:"reqId,omitempty" json:"reqId,omitempty"`
	Package          primitive.ObjectID    `bson:"package,omitempty" json:"package,omitempty"`
	Program          primitive.ObjectID    `bson:"program,omitempty" json:"program,omitempty"`
	InvoiceId        primitive.ObjectID    `bson:"invoiceId,omitempty" json:"invoiceId,omitempty"`
	Date             time.Time             `bson:"date,omitempty" json:"date,omitempty"`
	CustomerId       primitive.ObjectID    `bson:"customerId,omitempty" json:"customerId,omitempty"` //same as user id
	Status           enums.TravelReqStatus `bson:"status,omitempty" json:"status,omitempty"`
	TravelRequestDto `bson:",inline"`
	RevisionNum      int                `bson:"revisionNum,omitempty" json:"revisionNum,omitempty"`
	Trash            bool               `bson:"trash" json:"trash"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy        primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy        primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type TravelRequestPagination struct {
	Requests   []TravelRequest  `json:"requests"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

// change status
type ChangeStatusDto struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"`
}

func (c *ChangeStatusDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, c)
}

//

type TravelRequestWithProgram struct {
	TravelRequest `bson:",inline"`
	Program       Program               `bson:"program" json:"program"`
	Customer      membermodels.Customer `bson:"customer" json:"customer"`
}
