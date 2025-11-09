package enums

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
)

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "waitingforpayment"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusUnpaid  InvoiceStatus = "unpaid"
)