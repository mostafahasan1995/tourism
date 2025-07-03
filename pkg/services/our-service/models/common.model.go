package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Destination struct {
	DestinationFrom primitive.ObjectID `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo   primitive.ObjectID `bson:"destinationTo" json:"destinationTo"`
	TripDetails     TripDetails        `bson:"tripDetails" json:"tripDetails"`
	Accommodation   []Accommodation    `bson:"accommodation" json:"accommodation"`
	FlightTickets   FlightTicket       `bson:"flightTickets" json:"flightTickets"`
	Transportation  Transportation     `bson:"transportation" json:"transportation"`
	Activities      []string           `bson:"activities" json:"activities"`
	Agenda          Agenda             `bson:"agenda" json:"agenda"`
	Services        Services           `bson:"services" json:"services"`
}

type TripDetails struct {
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	Days      int       `bson:"days" json:"days"`
}

type Accommodation struct {
	City            string   `bson:"city" json:"city"`
	Preferences     string   `bson:"preferences" json:"preferences"`
	SpecialRequests []string `bson:"specialRequests" json:"specialRequests"`
	NumOfRooms      int      `bson:"numOfRooms" json:"numOfRooms"`
	BedType         string   `bson:"bedType" json:"bedType"`
	NumOfBeds       int      `bson:"numOfBeds" json:"numOfBeds"`
	NumOfBathrooms  int      `bson:"numOfBathrooms" json:"numOfBathrooms"`
	//delegation
	PrivateMeetingRoom bool `bson:"privateMeetingRoom" json:"privateMeetingRoom"`
}

type FlightTicket struct {
	ArrangeByUs            bool                       `bson:"arrangeByUs" json:"arrangeByUs"`
	DestinationFrom        primitive.ObjectID         `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo          primitive.ObjectID         `bson:"destinationTo" json:"destinationTo"`
	TripType               string                     `bson:"tripType" json:"tripType"`       // e.g round-trip - one-way - multi-city
	TravelClass            string                     `bson:"travelClass" json:"travelClass"` // e.g economy - business - first class
	FlexibleTravelDates    string                     `bson:"flexibleTravelDates" json:"flexibleTravelDates"`
	DepartureDate          *time.Time                 `bson:"departureDate" json:"departureDate"`
	ReturnDate             *time.Time                 `bson:"returnDate" json:"returnDate"`
	PreferredDepartureTime string                     `bson:"preferredDepartureTime" json:"preferredDepartureTime"` // e.g. evening - morning
	LayoverPreferences     string                     `bson:"layoverPreferences" json:"layoverPreferences"`         // short transit time
	PreferredAirlines      string                     `bson:"preferredAirlines" json:"preferredAirlines"`           // e.g. franch arilines
	ExtraBaggage           bool                       `bson:"extraBaggage" json:"extraBaggage"`
	SpecialMeals           bool                       `bson:"specialMeals" json:"specialMeals"`
	TravelWithPet          bool                       `bson:"travelWithPet" json:"travelWithPet"`
	Adults                 int                        `bson:"adults" json:"adults"`
	Children               int                        `bson:"children" json:"children"`
	Infant                 int                        `bson:"infant" json:"infant"`
	PreferredContactMethod string                     `bson:"preferredContactMethod" json:"preferredContactMethod"`
	SpecialRequirements    transl.Localizable[string] `bson:"specialRequirements" json:"specialRequirements"`
}

type Transportation struct {
	TransType             string             `bson:"transType" json:"transType"`
	Capacity              string             `bson:"capacity" json:"capacity"`
	DriverLanguagesSpoken []string           `bson:"driverLanguagesSpoken" json:"driverLanguagesSpoken"`
	LuxuryFeatures        []string           `bson:"luxuryFeatures" json:"luxuryFeatures"`
	StartDate             *time.Time         `bson:"startDate" json:"startDate"`
	EndDate               *time.Time         `bson:"endDate" json:"endDate"`
	StartTime             string             `bson:"startTime" json:"startTime"`
	EndTime               string             `bson:"endTime" json:"endTime"`
	CarTypeId             primitive.ObjectID `bson:"carTypeId" json:"carTypeId"`
}

type Agenda struct {
	NeedAgenda bool              `bson:"needAgenda" json:"needAgenda"`
	Items      []string          `bson:"items" json:"items"`
	Files      []types.FileField `bson:"files" json:"files"`
}

type Services struct {
	TourGuide           bool `bson:"tourGuide" json:"tourGuide"`
	Translator          bool `bson:"translator" json:"translator"`
	AirportPickup       bool `bson:"airportPickup" json:"airportPickup"`
	TourAfterMeeting    bool `bson:"tourAfterMeeting" json:"tourAfterMeeting"`
	Photography         bool `bson:"photography" json:"photography"`
	AirportMeetAndGreet bool `bson:"airportMeetAndGreet" json:"airportMeetAndGreet"`
	SimCardAndInternet  bool `bson:"simCardAndInternet" json:"simCardAndInternet"`
}

//

type Delegation struct {
	DelegationType      string                     `bson:"delegationType" json:"delegationType"`
	OrganizationName    transl.Localizable[string] `bson:"organizationName" json:"organizationName"`
	TripCoordinatorName string                     `bson:"tripCoordinatorName" json:"tripCoordinatorName"`
}

type BusinessMan struct {
	Purpose             string            `bson:"purpose" json:"purpose"` //meeting - investment - conference - other
	ClientIsCoordinator bool              `bson:"clientIsCoordinator" json:"clientIsCoordinator"`
	CoordinatorName     string            `bson:"coordinatorName" json:"coordinatorName"`
	CoordinatorPhone    types.PhoneNumber `bson:"coordinatorPhone" json:"coordinatorPhone"`
	CoordinatorEmail    string            `bson:"coordinatorEmail" json:"coordinatorEmail"`
}

type HotelBooking struct {
	HotelId   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
	HotelName string             `bson:"hotelName" json:"hotelName"`
}

type VipCar struct {
	Destinations []VipCarDest `bson:"destinations" json:"destinations"`
}

type VipCarDest struct {
	DestinationFrom primitive.ObjectID `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo   primitive.ObjectID `bson:"destinationTo" json:"destinationTo"`
	Transportation  `bson:",inline"`
	Services        Services `bson:"services" json:"services"`
}

// custom plan
type CustomPlan struct {
	TripType        string             `bson:"tripType" json:"tripType"` //e.g. family - honeymoon - luxury retreat - other
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator"`
}

// flight request
type FlightTicketRequest struct {
	Destinations []FlightTicketDestination `bson:"destinations" json:"destinations"`
}

type FlightTicketDestination struct {
	FlightTicket `bson:",inline"`
	Services     Services `bson:"services" json:"services"`
}

// partner request
type PartnerRequest struct {
	CompanyName     string        `bson:"companyName" json:"companyName"`
	BusinessType    string        `bson:"businessType" json:"businessType"`
	Website         string        `bson:"website" json:"website"`
	CompanyLocation string        `bson:"companyLocation" json:"companyLocation"`
	ContactDetail   ContactDetail `bson:"contactDetail" json:"contactDetail"`
	Message         string        `bson:"message" json:"message"`
}

type ContactDetail struct {
	FullName    string            `bson:"fullName" json:"fullName"`
	Position    string            `bson:"position" json:"position"`
	PhoneNumber types.PhoneNumber `bson:"phoneNumber" json:"phoneNumber"`
	Email       string            `bson:"email" json:"email"`
}

//
