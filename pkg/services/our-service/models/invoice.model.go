package models

import (
	"errors"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceDto struct {
	DateOfIssue time.Time           `bson:"dateOfIssue" json:"dateOfIssue"`
	Customer    InvoiceContact      `bson:"customer" json:"customer"`
	Company     InvoiceContact      `bson:"company" json:"company"`
	ProgramName string              `bson:"programName" json:"programName"`
	TravelStart time.Time           `bson:"travelStart" json:"travelStart"`
	TravelEnd   time.Time           `bson:"travelEnd" json:"travelEnd"`
	Services    []InvoiceService    `bson:"services" json:"services"`
	Adjustments []InvoiceAdjustment `bson:"adjustments" json:"adjustments"`
	Note        string              `bson:"note" json:"note"`
}

type InvoiceAdjustment struct {
	Type    string  `bson:"type" json:"type"` //discount - surcharge
	Amount  float64 `bson:"amount" json:"amount"`
	Percent float64 `bson:"percent" json:"percent"`
}

type InvoiceContact struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Email   string `bson:"email" json:"email"`
	Website string `bson:"website" json:"website"`
}

type InvoiceService struct {
	Item  string  `bson:"item" json:"item"`
	Price float64 `bson:"price" json:"price"`
	Qty   int     `bson:"qty" json:"qty"`
	Total float64 `bson:"total" json:"total"`
}

type Payment struct {
	PaymentId  string `bson:"paymentId" json:"paymentId"`
	PaymentDto `bson:",inline"`
}

type PaymentDto struct {
	Date   time.Time `bson:"date" json:"date"`
	Method string    `bson:"method" json:"method"`
	Amount float64   `bson:"amount" json:"amount"`
	Unit   string    `bson:"unit" json:"unit"`
	Status string    `bson:"status" json:"status"` //paid - unpaid
	Note   string    `bson:"note" json:"note"`
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

	//todo: add adjustments

	i.SubTotal = subTotal
	i.Total = subTotal

	var paidAmount float64
	for _, payment := range i.Payments {
		if payment.Status == "paid" {
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
