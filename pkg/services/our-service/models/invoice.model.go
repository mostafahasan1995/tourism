package models

import (
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceDto struct {
	DateOfIssue time.Time           `bson:"dateOfIssue" json:"dateOfIssue" validate:"required"`
	Customer    InvoiceContact      `bson:"customer" json:"customer" validate:"required,dive"`
	Company     InvoiceContact      `bson:"company" json:"company" validate:"required,dive"`
	ProgramName string              `bson:"programName" json:"programName" validate:"required"`
	TravelStart time.Time           `bson:"travelStart" json:"travelStart" validate:"required"`
	TravelEnd   time.Time           `bson:"travelEnd" json:"travelEnd" validate:"required"`
	Services    []InvoiceService    `bson:"services" json:"services"`
	Adjustments []InvoiceAdjustment `bson:"adjustments" json:"adjustments"`
	Note        string              `bson:"note" json:"note"`
}

type InvoiceAdjustment struct {
	Type    string  `bson:"type" json:"type"` //addition - substruction
	Title   string  `bson:"title" json:"title"`
	Amount  float64 `bson:"amount" json:"amount"`
	Percent float64 `bson:"percent" json:"percent"`
}

type InvoiceContact struct {
	Name    string            `bson:"name" json:"name" validate:"required"`
	Address string            `bson:"address" json:"address" validate:"required"`
	Phone   types.PhoneNumber `bson:"phone" json:"phone" validate:"required"`
	Email   string            `bson:"email" json:"email" validate:"required"`
	Website string            `bson:"website" json:"website"`
}

type InvoiceService struct {
	Item  string  `bson:"item" json:"item"`
	Price float64 `bson:"price" json:"price"`
	Qty   int     `bson:"qty" json:"qty"`
	Total float64 `bson:"total" json:"total"`
}

func (i *InvoiceDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, i)
}

type Invoice struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	InvoiceId    string             `bson:"invoiceId" json:"invoiceId"`
	InvoiceDto   `bson:",inline"`
	TravelReqId  primitive.ObjectID `bson:"travelReqId" json:"travelReqId"`
	SubTotal     float64            `bson:"subTotal" json:"subTotal"`
	Total        float64            `bson:"total" json:"total"`
	PaidAmount   float64            `bson:"paidAmount" json:"paidAmount"`
	UnpaidAmount float64            `bson:"unpaidAmount" json:"unpaidAmount"`
	Payments     []Payment          `bson:"payments" json:"payments"`
	Trash        bool               `bson:"trash" json:"trash"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy    primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt    time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy    primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

func (i *Invoice) SetTotals() error {
	var services []InvoiceService
	var subTotal float64
	for _, service := range i.Services {
		service.Total = service.Price * float64(service.Qty)
		services = append(services, service)
		subTotal += service.Total
	}

	i.Services = services

	total := subTotal

	for _, adjustment := range i.Adjustments {

		if adjustment.Type == "addition" {
			if adjustment.Percent > 0 {
				total += subTotal * adjustment.Percent / 100
			} else if adjustment.Amount > 0 {
				total += adjustment.Amount
			}
		} else if adjustment.Type == "substruction" {
			if adjustment.Percent > 0 {
				total -= subTotal * adjustment.Percent / 100
			} else if adjustment.Amount > 0 {
				total -= adjustment.Amount
			}
		}
	}

	i.SubTotal = subTotal
	i.Total = total

	var paidAmount float64
	for _, payment := range i.Payments {
		if payment.Status == enums.PaymentStatusPaid {
			paidAmount += payment.Amount
		}
	}

	if paidAmount > i.Total {
		return errors.New("paid amount is greater than total")
	}

	unpaidAmount := i.Total - paidAmount
	i.PaidAmount = paidAmount
	i.UnpaidAmount = unpaidAmount

	return nil
}

type InvoicePagination struct {
	Invoices   []Invoice        `bson:"invoices" json:"invoices"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

// payments
type Payment struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PaymentId  string             `bson:"paymentId" json:"paymentId"`
	PaymentDto `bson:",inline"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type PaymentDto struct {
	Date    time.Time           `bson:"date" json:"date" validate:"required"`
	Method  string              `bson:"method" json:"method" validate:"required"`
	Amount  float64             `bson:"amount" json:"amount" validate:"required"`
	Unit    string              `bson:"unit" json:"unit"`
	Status  enums.PaymentStatus `bson:"status" json:"status" validate:"required,oneof=paid unpaid"` //paid - unpaid
	Note    string              `bson:"note" json:"note"`
	Receipt []types.FileField   `bson:"receipt" json:"receipt"`
}

func (p *PaymentDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, p)
}

type PayOrder struct {
	Amount  float64           `bson:"amount" json:"amount"`
	Method  string            `bson:"method" json:"method"`
	Receipt []types.FileField `bson:"receipt" json:"receipt"`
}

//travel request agent

// type InvoiceTravelReqData struct {
// 	TravelReqId      primitive.ObjectID `bson:"travelReqId" json:"travelReqId"`
// 	DepartureAgent   primitive.ObjectID `bson:"departureAgent" json:"departureAgent"`
// 	DestinationAgent primitive.ObjectID `bson:"destinationAgent" json:"destinationAgent"`
// }
