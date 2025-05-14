package models

import (
	// "larsa-tourism-microservices/pkg/types"
	// "larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightTicketRequestDto struct {
	DestinationFrom string `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo   string `bson:"destinationTo" json:"destinationTo"`
	TripType        string `bson:"tripType" json:"tripType"`
	TravelClass     string `bson:"travelClass" json:"travelClass"`

	DepartureDate time.Time `bson:"departureDate" json:"departureDate"`
	ReturnDate    time.Time `bson:"returnDate" json:"returnDate"`

	NumberOfAdults   int `bson:"numberOfAdults" json:"numberOfAdults"`
	NumberOfChildren int `bson:"numberOfChildren" json:"numberOfChildren"`
	NumberOfInfants  int `bson:"numberOfInfants" json:"numberOfInfants"`

	BestDepartureTime      string `bson:"bestDepartureTime" json:"bestDepartureTime"`
	StopoverPreferences    string `bson:"stopoverPreferences" json:"stopoverPreferences"`
	PreferredContactMethod string `bson:"preferredContactMethod" json:"preferredContactMethod"`

	PreferredAirlines string `bson:"preferredAirlines" json:"preferredAirlines"`
	ExtraLuggage      bool   `bson:"extraLuggage" json:"extraLuggage"`
	SpecialMeals      bool   `bson:"specialMeals" json:"specialMeals"`
}

type FlightTicketRequest struct {
	FlightTicketRequestDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type FlightTicketRequestPagination struct {
	FlightTicketRequest []FlightTicketRequest `bson:"flightTicketRequest" json:"flightTicketRequest"`

	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
