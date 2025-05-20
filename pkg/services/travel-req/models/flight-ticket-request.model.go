package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightTicketRequestDto struct {
	Destinations []FlightTicktDest `bson:"destinations" json:"destinations"`
}

type FlightTicktDest struct {
	DestinationFrom primitive.ObjectID `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo   primitive.ObjectID `bson:"destinationTo" json:"destinationTo"`
	TripType        string             `bson:"tripType" json:"tripType"`
	TravelClass     string             `bson:"travelClass" json:"travelClass"`

	DepartureDate time.Time `bson:"departureDate" json:"departureDate"`
	ReturnDate    time.Time `bson:"returnDate" json:"returnDate"`

	NumberOfAdults   int `bson:"numberOfAdults" json:"numberOfAdults"`
	NumberOfChildren int `bson:"numberOfChildren" json:"numberOfChildren"`
	NumberOfInfants  int `bson:"numberOfInfants" json:"numberOfInfants"`

	BestDepartureTime   string `bson:"bestDepartureTime" json:"bestDepartureTime"`     //morning - afternoon - evening
	StopoverPreferences string `bson:"stopoverPreferences" json:"stopoverPreferences"` //shortest - cheapest - fastest
	PreferredAirlines   string `bson:"preferredAirlines" json:"preferredAirlines"`

	ExtraLuggage           bool   `bson:"extraLuggage" json:"extraLuggage"`
	SpecialMeals           bool   `bson:"specialMeals" json:"specialMeals"`
	PreferredContactMethod string `bson:"preferredContactMethod" json:"preferredContactMethod"`
}

type FlightTicketRequest struct {
	Id                     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FlightTicketRequestDto `bson:",inline"`
}
