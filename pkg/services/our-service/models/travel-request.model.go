package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Example JSON for TravelRequestDto:
// {
//   "clientName": "John Doe",
//   "clientPhone": "+1234567890",
//   "clientEmail": "john.doe@example.com",
//   "serviceType": "delegation",
//   "delegation": {
//     "delegationType": "business",
//     "organizationName": "Tech Corp",
//     "nationality": "USA",
//     "tripDuration": 5,
//     "tripCoordinatorName": "Jane Smith"
//   },
//   "destination": {
//     "destination": "507f1f77bcf86cd799439011",
//     "tripDetails": {
//       "startDate": "2024-03-01T00:00:00Z",
//       "endDate": "2024-03-05T00:00:00Z",
//       "days": 5
//     },
//     "accommodation": [
//       {
//         "city": "Dubai",
//         "preferences": "5-star hotel",
//         "specialRequests": ["Non-smoking room", "High floor"],
//         "numOfRooms": 2,
//         "bedType": "King",
//         "numOfBeds": 2,
//         "numOfBathrooms": 2,
//         "privateMeetingRoom": true
//       }
//     ],
//     "flightTickets": {
//       "arrangeByUs": true,
//       "tripType": "round-trip",
//       "travelClass": "business",
//       "numOfPassengers": 4,
//       "departureDate": "2024-03-01T00:00:00Z",
//       "returnDate": "2024-03-05T00:00:00Z",
//       "flexibleTravelDates": "No",
//       "preferredDepartureTime": "morning",
//       "layoverPreferences": "short transit time",
//       "preferredAirlines": "Emirates",
//       "extraBaggage": true,
//       "specialMeals": true,
//       "travelWithPet": false,
//       "anySpecialReq": "Vegetarian meals"
//     },
//     "transportation": {
//       "trans": "private car",
//       "driverLangs": ["English", "Arabic"]
//     },
//     "activities": ["Business meetings", "City tour"],
//     "agenda": {
//       "needAgenda": true,
//       "items": ["Morning meetings", "Afternoon site visits"],
//       "files": []
//     },
//     "services": {
//       "tourGuide": true,
//       "translator": true,
//       "airportPickup": true,
//       "tourAfterMeeting": true
//     }
//   },
//   "tripCoordinator": "507f1f77bcf86cd799439012",
//   "contactMethod": ["email", "phone"],
//   "specialReq": "Need wheelchair assistance at airport"
// }
//
// Note: This is an example for a delegation-type request. For other service types:
// - For businessMan: include businessMan object instead of delegation
// - For vipCar: include vipCar object with destinations
// - For customPlan: include customPlan object
// - For flightTicketRequest: include flightTicketRequest object
// - For partnerRequest: include partnerRequest object

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
