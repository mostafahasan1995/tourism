package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusinessFrom struct {
	Id                  primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Purpose             string               `bson:"purpose" json:"purpose"` //could be an ObjectId
	ClientName          string               `bson:"clientName" json:"clientName"`
	ClientPhone         string               `bson:"clientPhone" json:"clientPhone"`
	ClientEmail         string               `bson:"clientEmail" json:"clientEmail"`
	ClientIsCoordinator bool                 `bson:"clientIsCoordinator" json:"clientIsCoordinator"`
	CoordinatorName     string               `bson:"coordinatorName" json:"coordinatorName"`
	CoordinatorPhone    string               `bson:"coordinatorPhone" json:"coordinatorPhone"`
	CoordinatorEmail    string               `bson:"coordinatorEmail" json:"coordinatorEmail"`
	Nationality         string               `bson:"nationality" json:"nationality"`
	TripDuration        int                  `bson:"tripDuration" json:"tripDuration"`
	Destinations        []BFDestination      `bson:"destinations" json:"destinations"`
	Coordinator         []primitive.ObjectID `bson:"coordinator" json:"coordinator"`
	ContactMethod       []string             `bson:"contactMethod" json:"contactMethod"`
	SpecialReq          string               `bson:"specialReq" json:"specialReq"`
}

type BFDestination struct {
	Name           string         `bson:"name" json:"name"`
	TripDetails    TripDetails    `bson:"tripDetails" json:"tripDetails"`
	Accommodation  Accommodation  `bson:"accommodation" json:"accommodation"`
	FlightTickets  FlightTicket   `bson:"flightTickets" json:"flightTickets"`
	Transportation Transportation `bson:"transportation" json:"transportation"`
	Activities     []string       `bson:"activities" json:"activities"`
	Agenda         Agenda         `bson:"agenda" json:"agenda"`
	Services       Services       `bson:"services" json:"services"`
}

type TripDetails struct {
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	Days      int       `bson:"days" json:"days"`
}

type Accommodation struct {
	Preferences     string   `bson:"preferences" json:"preferences"`
	SpecialRequests []string `bson:"specialRequests" json:"specialRequests"`
	NumOfRooms      int      `bson:"numOfRooms" json:"numOfRooms"`
	BedType         string   `bson:"bedType" json:"bedType"`
	NumOfBeds       int      `bson:"numOfBeds" json:"numOfBeds"`
	NumOfBathrooms  int      `bson:"numOfBathrooms" json:"numOfBathrooms"`
}

type FlightTicket struct {
	ArrangeByUs         bool      `bson:"arrangeByUs" json:"arrangeByUs"`
	TripType            string    `bson:"tripType" json:"tripType"`
	TravelClass         string    `bson:"travelClass" json:"travelClass"`
	DepartureDate       time.Time `bson:"departureDate" json:"departureDate"`
	ReturnDate          time.Time `bson:"returnDate" json:"returnDate"`
	FlexibleTravelDates string    `bson:"flexibleTravelDates" json:"flexibleTravelDates"`
	//number of passengers
	NumberOfAdults   int `bson:"numberOfAdults" json:"numberOfAdults"`
	NumberOfChildren int `bson:"numberOfChildren" json:"numberOfChildren"`
	NumberOfInfants  int `bson:"numberOfInfants" json:"numberOfInfants"`
	//preferrd flight options
	DepartureTime string `bson:"departureTime" json:"departureTime"`
	Layover       string `bson:"layover" json:"layover"`
	Airlines      string `bson:"airlines" json:"airlines"`
	//
	ExtraBaggage  bool   `bson:"extraBaggage" json:"extraBaggage"`
	SpecialMeals  bool   `bson:"specialMeals" json:"specialMeals"`
	WithPet       bool   `bson:"withPet" json:"withPet"`
	AnySpecialReq string `bson:"anySpecialReq" json:"anySpecialReq"`
}

type Transportation struct {
	Trans       string   `bson:"trans" json:"trans"`
	DriverLangs []string `bson:"driverLangs" json:"driverLangs"`
}

type Agenda struct {
	NeedAgenda bool              `bson:"needAgenda" json:"needAgenda"`
	AgendaText []string          `bson:"agendaText" json:"agendaText"`
	Files      []types.FileField `bson:"files" json:"files"`
}

type Services struct {
	TourGuide        bool `bson:"tourGuide" json:"tourGuide"`
	Translator       bool `bson:"translator" json:"translator"`
	AirportPickup    bool `bson:"airportPickup" json:"airportPickup"`
	TourAfterMeeting bool `bson:"tourAfterMeeting" json:"tourAfterMeeting"`
}
