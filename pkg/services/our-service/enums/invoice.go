package enums

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
)

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "waiting-payment"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusUnpaid  InvoiceStatus = "unpaid"

	InvoiceStatusWaitingApproved InvoiceStatus = "waiting-approved"
)

type WithdrawalStatus string

const (
	WithdrawalStatusPending  WithdrawalStatus = "pending"
	WithdrawalStatusApproved WithdrawalStatus = "approved"
	WithdrawalStatusRejected WithdrawalStatus = "rejected"
)
