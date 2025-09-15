package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PreBookData struct {
	Data       PreBook `json:"data" bson:"data"`
	GuestLevel int     `json:"guestLevel" bson:"guestLevel"`
	Sandbox    bool    `json:"sandbox" bson:"sandbox"`
}

type PreBook struct {
	BoardChanged           bool       `json:"boardChanged" bson:"boardChanged"`
	CancellationChanged    bool       `json:"cancellationChanged" bson:"cancellationChanged"`
	Commission             float64    `json:"commission" bson:"commission"`
	Currency               string     `json:"currency" bson:"currency"`
	HotelID                string     `json:"hotelId" bson:"hotelId"`
	IsPackageRate          bool       `json:"isPackageRate" bson:"isPackageRate"`
	OfferID                string     `json:"offerId" bson:"offerId"`
	PaymentTypes           []string   `json:"paymentTypes" bson:"paymentTypes"`
	PrebookID              string     `json:"prebookId" bson:"prebookId"`
	Price                  float64    `json:"price" bson:"price"`
	PriceDifferencePercent float64    `json:"priceDifferencePercent" bson:"priceDifferencePercent"`
	PriceType              string     `json:"priceType" bson:"priceType"`
	RoomTypes              []RoomType `json:"roomTypes" bson:"roomTypes"`
	SecretKey              string     `json:"secretKey" bson:"secretKey"`
	SuggestedSellingPrice  float64    `json:"suggestedSellingPrice" bson:"suggestedSellingPrice"`
	Supplier               string     `json:"supplier" bson:"supplier"`
	SupplierID             int        `json:"supplierId" bson:"supplierId"`
	TermsAndConditions     string     `json:"termsAndConditions" bson:"termsAndConditions"`
	TransactionID          string     `json:"transactionId" bson:"transactionId"`
}

type UserPrebook struct {
	PreBook    `bson:",inline"`
	GusetLevel int                `bson:"gusetLevel" json:"gusetLevel"`
	UserId     primitive.ObjectID `bson:"userId" json:"userId"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	Trash      bool               `bson:"trash" json:"trash"`
}
