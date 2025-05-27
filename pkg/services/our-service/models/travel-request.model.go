package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelRequestDto struct {
	ClientInfo  `bson:",inline"`
	ServiceType string `bson:"serviceType" json:"serviceType"` // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	//
	Delegation          *Delegation          `bson:"delegation,omitempty" json:"delegation,omitempty"`
	BusinessMan         *BusinessMan         `bson:"businessMan,omitempty" json:"businessMan,omitempty"`
	VipCar              *VipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty"`
	CustomPlan          *CustomPlan          `bson:"customPlan,omitempty" json:"customPlan,omitempty"`
	FlightTicketRequest *FlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty"`
	PartnerRequest      *PartnerRequest      `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty"`
	//
	Destination     []Destination      `bson:"destination,omitempty" json:"destination,omitempty"`
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator"`
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}

/*
Example JSON structure for TravelRequestDto:
{
    "clientName": "John Doe",
    "clientPhone": "+1234567890",
    "clientEmail": "john.doe@example.com",
    "serviceType": "custom-plan",
    "delegation": null,
    "businessMan": null,
    "vipCar": null,
    "customPlan": {
        "tripType": "luxury",
        "nationality": "US",
        "tripDuration": 7,
        "tripCoordinator": "507f1f77bcf86cd799439011"
    },
    "flightTicketRequest": null,
    "partnerRequest": null,
    "destination": [
        {
            "destination": "507f1f77bcf86cd799439012",
            "tripDetails": {
                "startDate": "2024-06-01T00:00:00Z",
                "endDate": "2024-06-07T00:00:00Z",
                "days": 7
            },
            "accommodation": [
                {
                    "city": "Paris",
                    "preferences": "Luxury hotel",
                    "specialRequests": ["Ocean view", "King size bed"],
                    "numOfRooms": 1,
                    "bedType": "King",
                    "numOfBeds": 1,
                    "numOfBathrooms": 1,
                    "privateMeetingRoom": false
                }
            ],
            "flightTickets": {
                "arrangeByUs": true,
                "tripType": "round-trip",
                "travelClass": "business",
                "numOfPassengers": 2,
                "departureDate": "2024-06-01T08:00:00Z",
                "returnDate": "2024-06-07T16:00:00Z",
                "flexibleTravelDates": "No",
                "preferredDepartureTime": "morning",
                "layoverPreferences": "direct flight preferred",
                "preferredAirlines": "Emirates",
                "extraBaggage": true,
                "specialMeals": true,
                "travelWithPet": false,
                "anySpecialReq": "Window seats preferred",
                "adults": 2,
                "children": 0,
                "infant": 0
            },
            "transportation": {
                "trans": "Private car",
                "driverLangs": ["English", "French"]
            },
            "activities": ["City tour", "Museum visits", "Fine dining"],
            "agenda": {
                "needAgenda": true,
                "items": [
                    "Morning city tour",
                    "Afternoon museum visit",
                    "Evening fine dining"
                ],
                "files": [
                    {
                        "_id": "507f1f77bcf86cd799439013",
                        "originalName": "itinerary.pdf",
                        "path": "/uploads/itinerary.pdf",
                        "service": "file-service",
                        "expire": "2024-12-31T23:59:59Z",
                        "variants": []
                    }
                ]
            },
            "services": {
                "tourGuide": true,
                "translator": true,
                "airportPickup": true,
                "tourAfterMeeting": false,
                "photography": true,
                "airportMeetAndGreet": true,
                "simCardAndInternet": true
            }
        }
    ],
    "tripCoordinator": "507f1f77bcf86cd799439014",
    "contactMethod": ["email", "phone"],
    "specialReq": "Vegetarian meals preferred"
}
*/

type TravelRequest struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ReqId        string             `bson:"reqId" json:"reqId"`
	Package      primitive.ObjectID `bson:"package" json:"package"`
	Program      primitive.ObjectID `bson:"program" json:"program"`
	Date         time.Time          `bson:"date" json:"date"`
	CustomerName string             `bson:"customerName" json:"customerName"`
	CustomerId   primitive.ObjectID `bson:"customerId" json:"customerId"`
	Status       string             `bson:"status" json:"status"`

	TravelRequestDto `bson:",inline"`
}

type TravelRequestPagination struct {
	Requests   []TravelRequest  `json:"requests"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
