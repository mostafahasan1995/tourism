package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	membermodels "larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"errors"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelRequestDto struct {
	//basic information

	ClientName   string            `bson:"clientName" json:"clientName" validate:"required"`
	ClientPhone  types.PhoneNumber `bson:"clientPhone" json:"clientPhone" `
	ClientEmail  string            `bson:"clientEmail" json:"clientEmail" validate:"required"`
	Nationality  string            `bson:"nationality" json:"nationality" `
	TripDuration int               `bson:"tripDuration" json:"tripDuration" validate:"required"`
	//
	ServiceType enums.ServiceType `bson:"serviceType" json:"serviceType" validate:"required,oneof=delegation custom-plan business-man vip-car flight-request partner-request hotel-booking"` // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	//request
	Delegation   *Delegation   `bson:"delegation,omitempty" json:"delegation,omitempty" validate:"required_if=ServiceType delegation"`
	BusinessMan  *BusinessMan  `bson:"businessMan,omitempty" json:"businessMan,omitempty" validate:"required_if=ServiceType business-man"`
	CustomPlan   *CustomPlan   `bson:"customPlan,omitempty" json:"customPlan,omitempty" validate:"required_if=ServiceType custom-plan"`
	HotelBooking *HotelBooking `bson:"hotelBooking,omitempty" json:"hotelBooking,omitempty" validate:"required_if=ServiceType hotel-booking"`
	//
	VipCar              *VipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty" validate:"required_if=ServiceType vip-car"`
	FlightTicketRequest *FlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty" validate:"required_if=ServiceType flight-request"`
	PartnerRequest      *PartnerRequest      `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty" validate:"required_if=ServiceType partner-request"`
	//delegation - business-man - custom-plan info
	Destinations    []Destination      `bson:"destinations,omitempty" json:"destinations,omitempty" `
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator" validate:"required"` // destination agent id
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod" validate:"required"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}

func (t *TravelRequestDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, t)
}

type TravelRequest struct {
	Id               primitive.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	ReqId            string                `bson:"reqId" json:"reqId"`
	Package          primitive.ObjectID    `bson:"package" json:"package"`
	Program          primitive.ObjectID    `bson:"program" json:"program"`
	InvoiceId        primitive.ObjectID    `bson:"invoiceId" json:"invoiceId"`
	Date             time.Time             `bson:"date" json:"date"`
	CustomerId       primitive.ObjectID    `bson:"customerId" json:"customerId"` //same as user id
	DepartureAgent   primitive.ObjectID    `bson:"departureAgent" json:"departureAgent"`
	Status           enums.TravelReqStatus `bson:"status" json:"status"`
	TravelRequestDto `bson:",inline"`
	RejectReason     string             `bson:"rejectReason" json:"rejectReason"`
	RevisionNum      int                `bson:"revisionNum" json:"revisionNum"`
	Trash            bool               `bson:"trash" json:"trash"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy        primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy        primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

func (t *TravelRequest) GetDepartureDestinationId() (*primitive.ObjectID, error) {
	switch t.ServiceType {
	case enums.ServiceTypeDelegation, enums.ServiceTypeCustomPlan, enums.ServiceTypeBusinessMan, enums.ServiceTypeHotelBooking:
		if len(t.Destinations) == 0 {
			return nil, errors.New("no destinations found")
		}
		return &t.Destinations[0].DestinationFrom, nil

	case enums.ServiceTypeVipCar:
		if t.VipCar == nil {
			return nil, errors.New("vip car is nil")
		}
		if len(t.VipCar.Destinations) == 0 {
			return nil, errors.New("no vip car destinations found")
		}
		return &t.VipCar.Destinations[0].DestinationFrom, nil
	case enums.ServiceTypeFlightRequest:
		if t.FlightTicketRequest == nil {
			return nil, errors.New("flight ticket request is nil")
		}
		if len(t.FlightTicketRequest.Destinations) == 0 {
			return nil, errors.New("no flight ticket request destinations found")
		}
		return &t.FlightTicketRequest.Destinations[0].DestinationFrom, nil
	}

	return nil, errors.New("unsupported service type")
}

type TravelRequestRes struct {
	TravelRequest `bson:",inline"`
	ProgramData   Program               `bson:"programData" json:"programData"`
	CustomerData  membermodels.Customer `bson:"customerData" json:"customerData"`
	PackageData   Package               `bson:"packageData" json:"packageData"`
}

type TravelRequestPagination struct {
	Requests   []TravelRequestRes `json:"requests"`
	Pagination types.Pagination   `bson:"pagination" json:"pagination"`
}

// change status
type ChangeStatusDto struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"`
}

func (c *ChangeStatusDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, c)
}

type TravelRequestWithProgram struct {
	TravelRequest `bson:",inline"`
	Program       Program               `bson:"program" json:"program"`
	Customer      membermodels.Customer `bson:"customer" json:"customer"`
}

type RejectMyReq struct {
	Reason string `json:"reason"`
}

// customer requests
type CustomerTravelRequest struct {
	TravelRequestRes `bson:",inline"`
	Price            float64 `bson:"price" json:"price"`
}

type CustomerTravelRequestPagination struct {
	Requests   []CustomerTravelRequest `json:"requests"`
	Pagination types.Pagination        `bson:"pagination" json:"pagination"`
}

// agnet transactions

type AgentTransaction struct {
	TravelRequestId primitive.ObjectID `bson:"travelRequestId" json:"travelRequestId"`
	InvoiceId       primitive.ObjectID `bson:"invoiceId" json:"invoiceId"`
	Date            time.Time          `bson:"date" json:"date"`
	OrderId         string             `bson:"orderId" json:"orderId"`
	CustomerName    string             `bson:"customerName" json:"customerName"`
	Commission      float64            `bson:"commission" json:"commission"`
}

type AgentTransactionPagination struct {
	Transactions []AgentTransaction `bson:"transactions" json:"transactions"`
	Pagination   types.Pagination   `bson:"pagination" json:"pagination"`
}
