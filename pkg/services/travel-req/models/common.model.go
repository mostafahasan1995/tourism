package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CFDestination struct {
	Destination    primitive.ObjectID `bson:"destination" json:"destination"`
	TripDetails    TripDetails        `bson:"tripDetails" json:"tripDetails"`
	Accommodation  []Accommodation    `bson:"accommodation" json:"accommodation"`
	FlightTickets  FlightTicket       `bson:"flightTickets" json:"flightTickets"`
	Transportation Transportation     `bson:"transportation" json:"transportation"`
	Activities     []string           `bson:"activities" json:"activities"`
	Agenda         Agenda             `bson:"agenda" json:"agenda"`
	Services       Services           `bson:"services" json:"services"`
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
	ArrangeByUs bool   `bson:"arrangeByUs" json:"arrangeByUs"`
	TripType    string `bson:"tripType" json:"tripType"`       // e.g round-trip - one-way - multi-city
	TravelClass string `bson:"travelClass" json:"travelClass"` // e.g economy - business - first class
	//delegation
	NumOfPassengers int `bson:"numOfPassengers" json:"numOfPassengers"`
	//
	DepartureDate          string `bson:"departureDate" json:"departureDate"`
	ReturnDate             string `bson:"returnDate" json:"returnDate"`
	FlexibleTravelDates    string    `bson:"flexibleTravelDates" json:"flexibleTravelDates"`
	PreferredDepartureTime string    `bson:"preferredDepartureTime" json:"preferredDepartureTime"` // e.g. evening - morning
	LayoverPreferences     string    `bson:"layoverPreferences" json:"layoverPreferences"`         // short transit time
	PreferredAirlines      string    `bson:"preferredAirlines" json:"preferredAirlines"`           // e.g. franch arilines
	ExtraBaggage           bool      `bson:"extraBaggage" json:"extraBaggage"`
	SpecialMeals           bool      `bson:"specialMeals" json:"specialMeals"`
	TravelWithPet          bool      `bson:"travelWithPet" json:"travelWithPet"`
	AnySpecialReq          string    `bson:"anySpecialReq" json:"anySpecialReq"`
	//business man
	Adults   int `bson:"adults" json:"adults"`
	Children int `bson:"children" json:"children"`
	Infant   int `bson:"infant" json:"infant"`
}

type Transportation struct {
	Trans       string   `bson:"trans" json:"trans"`
	DriverLangs []string `bson:"driverLangs" json:"driverLangs"`
}

type Agenda struct {
	NeedAgenda bool              `bson:"needAgenda" json:"needAgenda"`
	Items      []string          `bson:"items" json:"items"`
	Files      []types.FileField `bson:"files" json:"files"`
}

type Services struct {
	TourGuide     bool `bson:"tourGuide" json:"tourGuide"`
	Translator    bool `bson:"translator" json:"translator"`
	AirportPickup bool `bson:"airportPickup" json:"airportPickup"`
	//business man
	TourAfterMeeting bool `bson:"tourAfterMeeting" json:"tourAfterMeeting"`
	//custom plan
	Photography         bool `bson:"photography" json:"photography"`
	AirportMeetAndGreet bool `bson:"airportMeetAndGreet" json:"airportMeetAndGreet"`
	SimCardAndInternet  bool `bson:"simCardAndInternet" json:"simCardAndInternet"`
}
