package models

import (
	"time"
)

// RateResponse represents the main response structure for hotel rates
type Rates struct {
	HotelID   string     `json:"hotelId" bson:"hotelId"`
	RoomTypes []RoomType `json:"roomTypes" bson:"roomTypes"`
	ET        int        `json:"et" bson:"et"`
	ExpiresAt time.Time  `bson:"expiresAt" json:"-"`
}

// RoomType represents a room type with its offers and rates
type RoomType struct {
	RoomTypeID            string      `json:"roomTypeId" bson:"roomTypeId"`
	OfferID               string      `json:"offerId" bson:"offerId"`
	Supplier              string      `json:"supplier" bson:"supplier"`
	SupplierID            int         `json:"supplierId" bson:"supplierId"`
	Rates                 []Rate      `json:"rates" bson:"rates"`
	OfferRetailRate       PriceAmount `json:"offerRetailRate" bson:"offerRetailRate"`
	SuggestedSellingPrice PriceAmount `json:"suggestedSellingPrice" bson:"suggestedSellingPrice"`
	OfferInitialPrice     PriceAmount `json:"offerInitialPrice" bson:"offerInitialPrice"`
	PriceType             string      `json:"priceType" bson:"priceType"`
	RateType              string      `json:"rateType" bson:"rateType"`
	PaymentTypes          []string    `json:"paymentTypes" bson:"paymentTypes"`
}

// Rate represents individual rate details
type Rate struct {
	RateID               string               `json:"rateId" bson:"rateId"`
	OccupancyNumber      int                  `json:"occupancyNumber" bson:"occupancyNumber"`
	Name                 string               `json:"name" bson:"name"`
	MaxOccupancy         int                  `json:"maxOccupancy" bson:"maxOccupancy"`
	AdultCount           int                  `json:"adultCount" bson:"adultCount"`
	ChildCount           int                  `json:"childCount" bson:"childCount"`
	BoardType            string               `json:"boardType" bson:"boardType"`
	BoardName            string               `json:"boardName" bson:"boardName"`
	Remarks              string               `json:"remarks" bson:"remarks"`
	PriceType            string               `json:"priceType" bson:"priceType"`
	Commission           []PriceAmount        `json:"commission" bson:"commission"`
	RetailRate           RetailRate           `json:"retailRate" bson:"retailRate"`
	CancellationPolicies CancellationPolicies `json:"cancellationPolicies" bson:"cancellationPolicies"`
	PaymentTypes         []string             `json:"paymentTypes" bson:"paymentTypes"`
	ProviderCommission   PriceAmount          `json:"providerCommission" bson:"providerCommission"`
}

// PriceAmount represents a price with currency
type PriceAmount struct {
	Amount   float64 `json:"amount" bson:"amount"`
	Currency string  `json:"currency" bson:"currency"`
	Source   string  `json:"source,omitempty" bson:"source,omitempty"`
}

// RetailRate represents the retail rate structure
type RetailRate struct {
	Total                 []PriceAmount `json:"total" bson:"total"`
	SuggestedSellingPrice []PriceAmount `json:"suggestedSellingPrice" bson:"suggestedSellingPrice"`
	InitialPrice          []PriceAmount `json:"initialPrice" bson:"initialPrice"`
	TaxesAndFees          []TaxFee      `json:"taxesAndFees" bson:"taxesAndFees"`
}

// TaxFee represents tax and fee information
type TaxFee struct {
	Included    bool    `json:"included" bson:"included"`
	Description string  `json:"description" bson:"description"`
	Amount      float64 `json:"amount" bson:"amount"`
	Currency    string  `json:"currency" bson:"currency"`
}

// CancellationPolicies represents cancellation policy information
type CancellationPolicies struct {
	CancelPolicyInfos []CancelPolicyInfo `json:"cancelPolicyInfos" bson:"cancelPolicyInfos"`
	HotelRemarks      []string           `json:"hotelRemarks" bson:"hotelRemarks"`
	RefundableTag     string             `json:"refundableTag" bson:"refundableTag"`
}

// CancelPolicyInfo represents individual cancellation policy details
type CancelPolicyInfo struct {
	CancelTime string  `json:"cancelTime" bson:"cancelTime"`
	Amount     float64 `json:"amount" bson:"amount"`
	Currency   string  `json:"currency" bson:"currency"`
	Type       string  `json:"type" bson:"type"`
	Timezone   string  `json:"timezone" bson:"timezone"`
}

// for now we will use this struct to get the rates
type RatesList struct {
	Data       []map[string]any `json:"data"`
	GuestLevel int              `json:"guestLevel"`
	Sandbox    bool             `json:"sandbox"`
}
