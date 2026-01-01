package models

import (
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvoice_SetTotals(t *testing.T) {
	tests := []struct {
		name    string
		invoice Invoice
		want    struct {
			subTotal     float64
			total        float64
			paidAmount   float64
			unpaidAmount float64
			err          bool
		}
	}{
		{
			name: "basic calculation with services only",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
						{Item: "Service 2", Price: 50, Qty: 3},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     350, // (100 * 2) + (50 * 3)
				total:        350,
				paidAmount:   0,
				unpaidAmount: 350,
				err:          false,
			},
		},
		{
			name: "calculation with percentage addition adjustment",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
					},
					Adjustments: []InvoiceAdjustment{
						{Type: "addition", Title: "Tax", Percent: 10},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     200,
				total:        220, // 200 + (200 * 0.1)
				paidAmount:   0,
				unpaidAmount: 220,
				err:          false,
			},
		},
		{
			name: "calculation with fixed amount subtraction adjustment",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
					},
					Adjustments: []InvoiceAdjustment{
						{Type: "substruction", Title: "Discount", Amount: 50},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     200,
				total:        150, // 200 - 50
				paidAmount:   0,
				unpaidAmount: 150,
				err:          false,
			},
		},
		{
			name: "calculation with partial payment",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
					},
				},
				Payments: []Payment{
					{
						PaymentDto: PaymentDto{
							Amount: 100,
							Status: enums.PaymentStatusPaid,
						},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     200,
				total:        200,
				paidAmount:   100,
				unpaidAmount: 100,
				err:          false,
			},
		},
		{
			name: "error when paid amount exceeds total",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
					},
				},
				Payments: []Payment{
					{
						PaymentDto: PaymentDto{
							Amount: 300,
							Status: enums.PaymentStatusPaid,
						},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     200,
				total:        200,
				paidAmount:   300,
				unpaidAmount: 0,
				err:          true,
			},
		},
		{
			name: "calculation with multiple adjustments",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
					},
					Adjustments: []InvoiceAdjustment{
						{Type: "addition", Title: "Tax", Percent: 10},
						{Type: "substruction", Title: "Discount", Amount: 30},
						{Type: "addition", Title: "Fee", Amount: 15},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     200,
				total:        205, // 200 + (200 * 0.1) - 30 + 15 = 200 + 20 - 30 + 15
				paidAmount:   0,
				unpaidAmount: 205,
				err:          false,
			},
		},
		{
			name: "calculation with unpaid payments ignored",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 2},
					},
				},
				Payments: []Payment{
					{
						PaymentDto: PaymentDto{
							Amount: 100,
							Status: enums.PaymentStatusPaid,
						},
					},
					{
						PaymentDto: PaymentDto{
							Amount: 50,
							Status: enums.PaymentStatusUnpaid,
						},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     200,
				total:        200,
				paidAmount:   100, // Only paid payments count
				unpaidAmount: 100,
				err:          false,
			},
		},
		{
			name: "calculation with percentage subtraction adjustment",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 3},
					},
					Adjustments: []InvoiceAdjustment{
						{Type: "substruction", Title: "Discount", Percent: 15},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     300,
				total:        255, // 300 - (300 * 0.15)
				paidAmount:   0,
				unpaidAmount: 255,
				err:          false,
			},
		},
		{
			name: "calculation with fixed amount addition adjustment",
			invoice: Invoice{
				InvoiceDto: InvoiceDto{
					Services: []InvoiceService{
						{Item: "Service 1", Price: 100, Qty: 1},
					},
					Adjustments: []InvoiceAdjustment{
						{Type: "addition", Title: "Service Fee", Amount: 25},
					},
				},
			},
			want: struct {
				subTotal     float64
				total        float64
				paidAmount   float64
				unpaidAmount float64
				err          bool
			}{
				subTotal:     100,
				total:        125, // 100 + 25
				paidAmount:   0,
				unpaidAmount: 125,
				err:          false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.invoice.SetTotals(0) // Use 0 profitRatio for tests to match old behavior

			if tt.want.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.subTotal, tt.invoice.SubTotal)
			assert.Equal(t, tt.want.total, tt.invoice.Total)
			assert.Equal(t, tt.want.paidAmount, tt.invoice.PaidAmount)
			assert.Equal(t, tt.want.unpaidAmount, tt.invoice.UnpaidAmount)

			// Verify service totals are calculated correctly
			for i := range tt.invoice.Services {
				expectedTotal := tt.invoice.Services[i].Price * float64(tt.invoice.Services[i].Qty)
				assert.Equal(t, expectedTotal, tt.invoice.Services[i].Total)
			}
		})
	}
}

func TestInvoice_SetTotals_AdjustmentValues(t *testing.T) {
	t.Run("adjustment values are calculated correctly", func(t *testing.T) {
		invoice := Invoice{
			InvoiceDto: InvoiceDto{
				Services: []InvoiceService{
					{Item: "Service 1", Price: 100, Qty: 2}, // subtotal = 200
				},
				Adjustments: []InvoiceAdjustment{
					{Type: "addition", Title: "Tax", Percent: 10},         // Val should be 20
					{Type: "substruction", Title: "Discount", Amount: 50}, // Val should be 50
					{Type: "addition", Title: "Fee", Amount: 30},          // Val should be 30
				},
			},
		}

		err := invoice.SetTotals(0) // Use 0 profitRatio for tests to match old behavior
		assert.NoError(t, err)

		// Check that Val fields are set correctly
		assert.Equal(t, 20.0, invoice.Adjustments[0].Val, "percentage addition adjustment Val")
		assert.Equal(t, 50.0, invoice.Adjustments[1].Val, "amount subtraction adjustment Val")
		assert.Equal(t, 30.0, invoice.Adjustments[2].Val, "amount addition adjustment Val")

		// Total should be: 200 + 20 - 50 + 30 = 200
		assert.Equal(t, 200.0, invoice.Total)
	})
}
