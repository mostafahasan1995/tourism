package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceDto struct {
	DateOfIssue   time.Time        `bson:"dateOfIssue" json:"dateOfIssue"`
	Customer      InvoiceContact   `bson:"customer" json:"customer"`
	Company       InvoiceContact   `bson:"company" json:"company"`
	Services      []InvoiceService `bson:"services" json:"services"`
	PaymentMethod string           `bson:"paymentMethod" json:"paymentMethod"`
	PaymentStatus string           `bson:"paymentStatus" json:"paymentStatus"`
	PaymentDate   time.Time        `bson:"paymentDate" json:"paymentDate"`
}

type InvoiceContact struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Email   string `bson:"email" json:"email"`
	Website string `bson:"website" json:"website"`
}

type InvoiceService struct {
	Name      string    `bson:"name" json:"name"`
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	Total     float64   `bson:"total" json:"total"`
}

type Invoice struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	InvoiceId  string             `bson:"invoiceId" json:"invoiceId"`
	InvoiceDto `bson:",inline"`
}

type InvoicePagination struct {
	Invoices   []Invoice        `bson:"invoices" json:"invoices"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
