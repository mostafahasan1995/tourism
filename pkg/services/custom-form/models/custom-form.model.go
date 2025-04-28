package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomFormType string

const (
	DelegationType CustomFormType = "delegation"
	BusinessType   CustomFormType = "business"
	PlanType       CustomFormType = "plan"
)

type Delegation struct {
	DelegationType  string `bson:"delegationType" json:"delegationType"`
	Organization    string `bson:"organization" json:"organization"`
	CoordinatorName string `bson:"coordinatorName" json:"coordinatorName"`
	Phone           string `bson:"phone" json:"phone"`
	Email           string `bson:"email" json:"email"`
}

type Business struct {
	Purpose             string `bson:"purpose" json:"purpose"` //could be an ObjectId
	ClientName          string `bson:"clientName" json:"clientName"`
	ClientPhone         string `bson:"clientPhone" json:"clientPhone"`
	ClientEmail         string `bson:"clientEmail" json:"clientEmail"`
	ClientIsCoordinator bool   `bson:"clientIsCoordinator" json:"clientIsCoordinator"`
	CoordinatorName     string `bson:"coordinatorName" json:"coordinatorName"`
	CoordinatorPhone    string `bson:"coordinatorPhone" json:"coordinatorPhone"`
	CoordinatorEmail    string `bson:"coordinatorEmail" json:"coordinatorEmail"`
}

type Plan struct {
	TripType   string `bson:"tripType" json:"tripType"`
	ClientName string `bson:"clientName" json:"clientName"`
	Phone      string `bson:"phone" json:"phone"`
	Email      string `bson:"email" json:"email"`
}

type CustomFormDto struct {
	Type          CustomFormType       `bson:"type" json:"type"` //delegation - business - plan
	Delegation    Delegation           `bson:"delegation" json:"delegation"`
	Business      Business             `bson:"business" json:"business"`
	Plan          Plan                 `bson:"plan" json:"plan"`
	Nationality   string               `bson:"nationality" json:"nationality"`
	TripDuration  int                  `bson:"tripDuration" json:"tripDuration"`
	Destinations  []BFDestination      `bson:"destinations" json:"destinations"`
	Coordinator   []primitive.ObjectID `bson:"coordinator" json:"coordinator"`
	ContactMethod []string             `bson:"contactMethod" json:"contactMethod"`
	SpecialReq    string               `bson:"specialReq" json:"specialReq"`
}

type CustomForm struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CustomFormDto `bson:",inline"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
}

type BFDestination struct {
	Name           string          `bson:"name" json:"name"`
	TripDetails    TripDetails     `bson:"tripDetails" json:"tripDetails"`
	Accommodation  []Accommodation `bson:"accommodation" json:"accommodation"`
	FlightTickets  FlightTicket    `bson:"flightTickets" json:"flightTickets"`
	Transportation Transportation  `bson:"transportation" json:"transportation"`
	Activities     []string        `bson:"activities" json:"activities"`
	Agenda         Agenda          `bson:"agenda" json:"agenda"`
	Services       Services        `bson:"services" json:"services"`
}

type TripDetails struct {
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	Days      int       `bson:"days" json:"days"`
}

type Accommodation struct {
	City               string   `bson:"city" json:"city"`
	Preferences        string   `bson:"preferences" json:"preferences"`
	SpecialRequests    []string `bson:"specialRequests" json:"specialRequests"`
	NumOfRooms         int      `bson:"numOfRooms" json:"numOfRooms"`
	BedType            string   `bson:"bedType" json:"bedType"`
	NumOfBeds          int      `bson:"numOfBeds" json:"numOfBeds"`
	NumOfBathrooms     int      `bson:"numOfBathrooms" json:"numOfBathrooms"`
	PrivateMeetingRoom bool     `bson:"privateMeetingRoom" json:"privateMeetingRoom"`
}

type FlightTicket struct {
	ArrangeByUs         bool      `bson:"arrangeByUs" json:"arrangeByUs"`
	TripType            string    `bson:"tripType" json:"tripType"`
	TravelClass         string    `bson:"travelClass" json:"travelClass"`
	DepartureDate       time.Time `bson:"departureDate" json:"departureDate"`
	ReturnDate          time.Time `bson:"returnDate" json:"returnDate"`
	FlexibleTravelDates string    `bson:"flexibleTravelDates" json:"flexibleTravelDates"`
	//number of passengers
	NumOfPassengers int `bson:"numOfPassengers" json:"numOfPassengers"`
	NumOfAdults     int `bson:"numOfAdults" json:"numOfAdults"`
	NumOfChildren   int `bson:"numOfChildren" json:"numOfChildren"`
	NumOfInfants    int `bson:"numOfInfants" json:"numOfInfants"`
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
	Items      []string          `bson:"items" json:"items"`
	Files      []types.FileField `bson:"files" json:"files"`
}

type Services struct {
	TourGuide           bool `bson:"tourGuide" json:"tourGuide"`
	Translator          bool `bson:"translator" json:"translator"`
	AirportPickup       bool `bson:"airportPickup" json:"airportPickup"`
	TourAfterMeeting    bool `bson:"tourAfterMeeting" json:"tourAfterMeeting"`
	PhotoAndVideo       bool `bson:"photoAndVideo" json:"photoAndVideo"`
	AirportMeetAndGreet bool `bson:"airportMeetAndGreet" json:"airportMeetAndGreet"`
	SIMAndInternet      bool `bson:"simAndInternet" json:"simAndInternet"`
}
