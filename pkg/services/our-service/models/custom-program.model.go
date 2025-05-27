package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CustomProgram struct {
	Delegation          Delegation                 `bson:"delegation,omitempty" json:"delegation,omitempty"`
	BusinessMan         BusinessMan                `bson:"businessMan,omitempty" json:"businessMan,omitempty"`
	CustomPlan          CustomPlan                 `bson:"customPlan,omitempty" json:"customPlan,omitempty"`
	VipCar              ProgramVipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty"`
	FlightTicketRequest ProgramFlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty"`
	PartnerRequest      PartnerRequest             `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty"`
	Destinations        []ProgramDestination       `bson:"destinations,omitempty" json:"destinations,omitempty"`
}

type ProgramVipCar struct {
	Destinations []ProgramVipCarDestination `bson:"destinations" json:"destinations"`
}

type ProgramVipCarDestination struct {
	DestinationFrom       primitive.ObjectID `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo         primitive.ObjectID `bson:"destinationTo" json:"destinationTo"`
	ProgramTransportation `bson:",inline"`
}

type ProgramFlightTicketRequest struct {
	Destinations []ProgramFlightTicket `bson:"destinations" json:"destinations"`
}

type ProgramDestination struct {
	DestinationFrom primitive.ObjectID     `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo   primitive.ObjectID     `bson:"destinationTo" json:"destinationTo"`
	TripDetails     TripDetails            `bson:"tripDetails" json:"tripDetails"`
	Accommodation   []ProgramAccommodation `bson:"accommodation" json:"accommodation"`
	FlightTickets   ProgramFlightTicket    `bson:"flightTickets" json:"flightTickets"`
	Transportation  ProgramTransportation  `bson:"transportation" json:"transportation"`
	Activities      ProgramActivities      `bson:"activities" json:"activities"`
	Agenda          Agenda                 `bson:"agenda" json:"agenda"`
	Services        ProgramServices        `bson:"services" json:"services"`
}

type ProgramAccommodation struct {
	Accommodation `bson:",inline"`
	PricePerNight float64 `bson:"pricePerNight" json:"pricePerNight"`
	TotalStayCost float64 `bson:"totalStayCost" json:"totalStayCost"`
}

type ProgramFlightTicket struct {
	FlightTicket `bson:",inline"`
	TotalCost    float64 `bson:"totalCost" json:"totalCost"`
}

type ProgramTransportation struct {
	Transportation `bson:",inline"`
	TotalCost      float64 `bson:"totalCost" json:"totalCost"`
}

type ProgramActivities struct {
	Activities []string `bson:"activities" json:"activities"`
	TotalCost  float64  `bson:"totalCost" json:"totalCost"`
}

type Service struct {
	Active bool    `bson:"active" json:"active"`
	Cost   float64 `bson:"cost" json:"cost"`
}

type ProgramServices struct {
	TourGuide           Service `bson:"tourGuide" json:"tourGuide"`
	Translator          Service `bson:"translator" json:"translator"`
	AirportPickup       Service `bson:"airportPickup" json:"airportPickup"`
	TourAfterMeeting    Service `bson:"tourAfterMeeting" json:"tourAfterMeeting"`
	Photography         Service `bson:"photography" json:"photography"`
	AirportMeetAndGreet Service `bson:"airportMeetAndGreet" json:"airportMeetAndGreet"`
	SimCardAndInternet  Service `bson:"simCardAndInternet" json:"simCardAndInternet"`
}
