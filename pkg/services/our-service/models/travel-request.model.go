package models

import (
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelRequestDto struct {
	//basic information
	ClientName   string `bson:"clientName" json:"clientName" validate:"required"`
	ClientPhone  string `bson:"clientPhone" json:"clientPhone" `
	ClientEmail  string `bson:"clientEmail" json:"clientEmail" validate:"required, email"`
	Nationality  string `bson:"nationality" json:"nationality" `
	TripDuration int    `bson:"tripDuration" json:"tripDuration" validate:"required"`
	//
	ServiceType enums.ServiceType `bson:"serviceType" json:"serviceType" validate:"required,oneof=delegation custom-plan business-man vip-car flight-request partner-request"` // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	//request
	Delegation          *Delegation          `bson:"delegation,omitempty" json:"delegation,omitempty" validate:"required_if=ServiceType delegation"`
	BusinessMan         *BusinessMan         `bson:"businessMan,omitempty" json:"businessMan,omitempty" validate:"required_if=ServiceType business-man"`
	VipCar              *VipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty" validate:"required_if=ServiceType vip-car"`
	CustomPlan          *CustomPlan          `bson:"customPlan,omitempty" json:"customPlan,omitempty" validate:"required_if=ServiceType custom-plan"`
	FlightTicketRequest *FlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty" validate:"required_if=ServiceType flight-request"`
	PartnerRequest      *PartnerRequest      `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty" validate:"required_if=ServiceType partner-request"`
	//
	Destination     []Destination      `bson:"destination,omitempty" json:"destination,omitempty" `
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator" validate:"required"`
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod" validate:"required"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}
type TravelRequest struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ReqId        string             `bson:"reqId" json:"reqId"`
	Package      primitive.ObjectID `bson:"package" json:"package"`
	Program      primitive.ObjectID `bson:"program" json:"program"`
	Date         time.Time          `bson:"date" json:"date"`
	CustomerName string             `bson:"customerName" json:"customerName"`
	CustomerId   primitive.ObjectID `bson:"customerId" json:"customerId"`
	Status       string             `bson:"status" json:"status"`

	TravelRequestDto `bson:",inline"`
}

type TravelRequestPagination struct {
	Requests   []TravelRequest  `json:"requests"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
