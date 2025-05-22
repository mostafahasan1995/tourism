package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Example JSON:
// {
//   "destinations": [
//     {
//       "destinationFrom": "507f1f77bcf86cd799439011",
//       "destinationTo": "507f1f77bcf86cd799439012",
//       "tripType": "round-trip",
//       "travelClass": "business",
//       "departureDate": "2024-03-20T08:00:00Z",
//       "returnDate": "2024-03-25T16:00:00Z",
//       "numberOfAdults": 2,
//       "numberOfChildren": 1,
//       "numberOfInfants": 0,
//       "bestDepartureTime": "morning",
//       "stopoverPreferences": "shortest",
//       "preferredAirlines": "Emirates",
//       "extraLuggage": true,
//       "specialMeals": true,
//       "preferredContactMethod": "email"
//     },
//     {
//       "destinationFrom": "507f1f77bcf86cd799439012",
//       "destinationTo": "507f1f77bcf86cd799439013",
//       "tripType": "one-way",
//       "travelClass": "economy",
//       "departureDate": "2024-03-26T10:00:00Z",
//       "returnDate": "2024-03-26T12:00:00Z",
//       "numberOfAdults": 2,
//       "numberOfChildren": 1,
//       "numberOfInfants": 0,
//       "bestDepartureTime": "morning",
//       "stopoverPreferences": "cheapest",
//       "preferredAirlines": "Qatar Airways",
//       "extraLuggage": false,
//       "specialMeals": true,
//       "preferredContactMethod": "phone"
//     }
//   ]
// }

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
