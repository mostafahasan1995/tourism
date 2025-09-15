package models

import "time"

type BookingData struct {
	Data       Booking `json:"data" bson:"data"`
	GuestLevel int     `json:"guestLevel" bson:"guestLevel"`
}

// Booking represents a hotel booking with all its details
type Booking struct {
	Addons                   any                         `json:"addons" bson:"addons"`
	AddonsRedemptions        any                         `json:"addonsRedemptions" bson:"addonsRedemptions"`
	AddonsTotalAmount        float64                     `json:"addonsTotalAmount" bson:"addonsTotalAmount"`
	Adults                   int                         `json:"adults" bson:"adults"`
	AgentID                  any                         `json:"agentId" bson:"agentId"`
	AmountRefunded           float64                     `json:"amountRefunded" bson:"amountRefunded"`
	APICommission            float64                     `json:"apiCommission" bson:"apiCommission"`
	BookedRooms              []BookedRoom                `json:"bookedRooms" bson:"bookedRooms"`
	BookingID                string                      `json:"bookingId" bson:"bookingId"`
	CancellationPolicies     BookingCancellationPolicies `json:"cancellationPolicies" bson:"cancellationPolicies"`
	CancelledAt              *time.Time                  `json:"cancelledAt" bson:"cancelledAt"`
	CancelledBy              any                         `json:"cancelledBy" bson:"cancelledBy"`
	Checkin                  string                      `json:"checkin" bson:"checkin"`
	Checkout                 string                      `json:"checkout" bson:"checkout"`
	Children                 string                      `json:"children" bson:"children"`
	ChildrenCount            int                         `json:"childrenCount" bson:"childrenCount"`
	ClientCommission         float64                     `json:"clientCommission" bson:"clientCommission"`
	ClientReference          string                      `json:"clientReference" bson:"clientReference"`
	Commission               float64                     `json:"commission" bson:"commission"`
	CreatedAt                string                      `json:"createdAt" bson:"createdAt"`
	Currency                 string                      `json:"currency" bson:"currency"`
	DistributorCommission    float64                     `json:"distributorCommission" bson:"distributorCommission"`
	DistributorPrice         float64                     `json:"distributorPrice" bson:"distributorPrice"`
	Email                    string                      `json:"email" bson:"email"`
	ExchangeRate             float64                     `json:"exchangeRate" bson:"exchangeRate"`
	ExchangeRateUsd          float64                     `json:"exchangeRateUsd" bson:"exchangeRateUsd"`
	FirstName                string                      `json:"firstName" bson:"firstName"`
	GuestID                  int                         `json:"guestId" bson:"guestId"`
	Holder                   BookingHolder               `json:"holder" bson:"holder"`
	HolderTitle              string                      `json:"holderTitle" bson:"holderTitle"`
	Hotel                    BookingHotel                `json:"hotel" bson:"hotel"`
	HotelConfirmationCode    string                      `json:"hotelConfirmationCode" bson:"hotelConfirmationCode"`
	HotelID                  string                      `json:"hotelId" bson:"hotelId"`
	HotelName                string                      `json:"hotelName" bson:"hotelName"`
	KnowBeforeYouGo          string                      `json:"knowBeforeYouGo" bson:"knowBeforeYouGo"`
	LastFreeCancellationDate string                      `json:"lastFreeCancellationDate" bson:"lastFreeCancellationDate"`
	LastName                 string                      `json:"lastName" bson:"lastName"`
	LoyaltyGuestID           any                         `json:"loyaltyGuestId" bson:"loyaltyGuestId"`
	MandatoryFees            string                      `json:"mandatoryFees" bson:"mandatoryFees"`
	Nationality              string                      `json:"nationality" bson:"nationality"`
	OptionalFees             string                      `json:"optionalFees" bson:"optionalFees"`
	PaymentScheduledAt       any                         `json:"paymentScheduledAt" bson:"paymentScheduledAt"`
	PaymentStatus            string                      `json:"paymentStatus" bson:"paymentStatus"`
	PaymentTransactionID     string                      `json:"paymentTransactionId" bson:"paymentTransactionId"`
	PrebookID                string                      `json:"prebookId" bson:"prebookId"`
	Price                    float64                     `json:"price" bson:"price"`
	ProcessingFee            float64                     `json:"processingFee" bson:"processingFee"`
	RebookFrom               string                      `json:"rebookFrom" bson:"rebookFrom"`
	RefundType               string                      `json:"refundType" bson:"refundType"`
	RefundedAt               any                         `json:"refundedAt" bson:"refundedAt"`
	Remarks                  string                      `json:"remarks" bson:"remarks"`
	Sandbox                  int                         `json:"sandbox" bson:"sandbox"`
	SellingPrice             string                      `json:"sellingPrice" bson:"sellingPrice"`
	SpecialRemarks           string                      `json:"specialRemarks" bson:"specialRemarks"`
	Status                   string                      `json:"status" bson:"status"`
	Supplier                 string                      `json:"supplier" bson:"supplier"`
	SupplierBookingID        string                      `json:"supplierBookingId" bson:"supplierBookingId"`
	SupplierBookingName      string                      `json:"supplierBookingName" bson:"supplierBookingName"`
	SupplierID               int                         `json:"supplierId" bson:"supplierId"`
	Tag                      string                      `json:"tag" bson:"tag"`
	TrackingID               string                      `json:"trackingId" bson:"trackingId"`
	UpdatedAt                string                      `json:"updatedAt" bson:"updatedAt"`
	UserID                   int                         `json:"userId" bson:"userId"`
	VoucherCode              string                      `json:"voucherCode" bson:"voucherCode"`
	VoucherID                any                         `json:"voucherId" bson:"voucherId"`
	VoucherTotalAmount       float64                     `json:"voucherTotalAmount" bson:"voucherTotalAmount"`
	VoucherTransationID      any                         `json:"voucherTransationId" bson:"voucherTransationId"`
}

// BookedRoom represents a single room booking
type BookedRoom struct {
	Adults               int                         `json:"adults" bson:"adults"`
	Amount               float64                     `json:"amount" bson:"amount"`
	Board                string                      `json:"board" bson:"board"`
	BoardCode            string                      `json:"boardCode" bson:"boardCode"`
	BoardName            string                      `json:"boardName" bson:"boardName"`
	BoardType            string                      `json:"boardType" bson:"boardType"`
	CancellationPolicies BookingCancellationPolicies `json:"cancellationPolicies" bson:"cancellationPolicies"`
	Children             int                         `json:"children" bson:"children"`
	ChildrenAges         any                         `json:"childrenAges" bson:"childrenAges"`
	ChildrenCount        int                         `json:"children_count" bson:"children_count"`
	Currency             string                      `json:"currency" bson:"currency"`
	FirstName            string                      `json:"firstName" bson:"firstName"`
	Guests               []BookingGuest              `json:"guests" bson:"guests"`
	LastName             string                      `json:"lastName" bson:"lastName"`
	OccupancyNumber      int                         `json:"occupancy_number" bson:"occupancy_number"`
	Rate                 BookingRate                 `json:"rate" bson:"rate"`
	Remarks              string                      `json:"remarks" bson:"remarks"`
	RoomType             BookingRoomType             `json:"roomType" bson:"roomType"`
	RoomID               string                      `json:"room_id" bson:"room_id"`
}

// BookingCancellationPolicies represents the cancellation policy details
type BookingCancellationPolicies struct {
	CancelPolicyInfos []BookingCancelPolicyInfo `json:"cancelPolicyInfos" bson:"cancelPolicyInfos"`
	HotelRemarks      any                       `json:"hotelRemarks" bson:"hotelRemarks"`
	RefundableTag     string                    `json:"refundableTag" bson:"refundableTag"`
}

// BookingCancelPolicyInfo represents individual cancellation policy information
type BookingCancelPolicyInfo struct {
	Amount     float64 `json:"amount" bson:"amount"`
	CancelTime string  `json:"cancelTime" bson:"cancelTime"`
	Currency   string  `json:"currency" bson:"currency"`
	Timezone   string  `json:"timezone" bson:"timezone"`
	Type       string  `json:"type" bson:"type"`
}

// BookingGuest represents a guest's information
type BookingGuest struct {
	Email           string `json:"email" bson:"email"`
	FirstName       string `json:"firstName" bson:"firstName"`
	LastName        string `json:"lastName" bson:"lastName"`
	OccupancyNumber int    `json:"occupancyNumber" bson:"occupancyNumber"`
	Phone           string `json:"phone" bson:"phone"`
	Remarks         string `json:"remarks" bson:"remarks"`
}

// BookingRate represents rate information for a room
type BookingRate struct {
	BoardName            string                      `json:"boardName" bson:"boardName"`
	BoardType            string                      `json:"boardType" bson:"boardType"`
	CancellationPolicies BookingCancellationPolicies `json:"cancellationPolicies" bson:"cancellationPolicies"`
	MaxOccupancy         int                         `json:"maxOccupancy" bson:"maxOccupancy"`
	RateID               string                      `json:"rateId" bson:"rateId"`
	Remarks              string                      `json:"remarks" bson:"remarks"`
	RetailRate           BookingRetailRate           `json:"retailRate" bson:"retailRate"`
}

// BookingRetailRate represents retail rate information
type BookingRetailRate struct {
	SuggestedSellingPrice BookingSuggestedSellingPrice `json:"suggestedSellingPrice" bson:"suggestedSellingPrice"`
	Total                 BookingTotal                 `json:"total" bson:"total"`
}

// BookingSuggestedSellingPrice represents suggested selling price information
type BookingSuggestedSellingPrice struct {
	Source string `json:"source" bson:"source"`
}

// BookingTotal represents total price information
type BookingTotal struct {
	Amount   float64 `json:"amount" bson:"amount"`
	Currency string  `json:"currency" bson:"currency"`
}

// BookingRoomType represents room type information
type BookingRoomType struct {
	Name       string `json:"name" bson:"name"`
	RoomTypeID string `json:"roomTypeId" bson:"roomTypeId"`
}

// BookingHolder represents the booking holder's information
type BookingHolder struct {
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Phone     string `json:"phone" bson:"phone"`
}

// BookingHotel represents hotel information
type BookingHotel struct {
	HotelID string `json:"hotelId" bson:"hotelId"`
	Name    string `json:"name" bson:"name"`
}
