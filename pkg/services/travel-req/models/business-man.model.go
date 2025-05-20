package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Example JSON:
// {
//   "purpose": "meeting",
//   "clientName": "John Smith",
//   "clientPhone": "+1234567890",
//   "clientEmail": "john.smith@company.com",
//   "clientIsCoordinator": false,
//   "coordinatorName": "Sarah Johnson",
//   "coordinatorPhone": "+1987654321",
//   "coordinatorEmail": "sarah.j@company.com",
//   "nationality": "American",
//   "tripDuration": 5,
//   "destinations": [
//     {
//       "destination": "507f1f77bcf86cd799439011",
//       "tripDetails": {
//         "startDate": "2024-03-20T10:00:00Z",
//         "endDate": "2024-03-25T14:00:00Z",
//         "days": 5
//       },
//       "accommodation": [
//         {
//           "city": "New York",
//           "preferences": "Business hotel in downtown",
//           "specialRequests": ["High floor", "Quiet room"],
//           "numOfRooms": 1,
//           "bedType": "King",
//           "numOfBeds": 1,
//           "numOfBathrooms": 1,
//           "privateMeetingRoom": true
//         }
//       ],
//       "flightTickets": {
//         "arrangeByUs": true,
//         "tripType": "round-trip",
//         "travelClass": "business",
//         "numOfPassengers": 1,
//         "departureDate": "2024-03-20T08:00:00Z",
//         "returnDate": "2024-03-25T16:00:00Z",
//         "flexibleTravelDates": "No",
//         "preferredDepartureTime": "morning",
//         "layoverPreferences": "short transit time",
//         "preferredAirlines": "Emirates",
//         "extraBaggage": true,
//         "specialMeals": true,
//         "travelWithPet": false,
//         "anySpecialReq": "Window seat preferred",
//         "adults": 1,
//         "children": 0,
//         "infant": 0
//       },
//       "transportation": {
//         "trans": "Private car",
//         "driverLangs": ["English", "Arabic"]
//       },
//       "activities": ["City tour", "Business networking"],
//       "agenda": {
//         "needAgenda": true,
//         "items": ["Morning meetings", "Afternoon site visits"],
//         "files": []
//       },
//       "services": {
//         "tourGuide": true,
//         "translator": true,
//         "airportPickup": true,
//         "tourAfterMeeting": true,
//         "photography": false,
//         "airportMeetAndGreet": true,
//         "simCardAndInternet": true
//       }
//     }
//   ],
//   "tripCoordinator": "507f1f77bcf86cd799439012",
//   "contactMethod": ["email", "phone"],
//   "specialReq": "Need vegetarian meals and airport transfer"
// }

type BusinessManDto struct {
	Purpose             string `bson:"purpose" json:"purpose"` //meeting - investment - conference - other
	ClientName          string `bson:"clientName" json:"clientName"`
	ClientPhone         string `bson:"clientPhone" json:"clientPhone"`
	ClientEmail         string `bson:"clientEmail" json:"clientEmail"`
	ClientIsCoordinator bool   `bson:"clientIsCoordinator" json:"clientIsCoordinator"`
	CoordinatorName     string `bson:"coordinatorName" json:"coordinatorName"`
	CoordinatorPhone    string `bson:"coordinatorPhone" json:"coordinatorPhone"`
	CoordinatorEmail    string `bson:"coordinatorEmail" json:"coordinatorEmail"`
	Nationality         string `bson:"nationality" json:"nationality"`
	TripDuration        int    `bson:"tripDuration" json:"tripDuration"`

	Destinations    []CFDestination    `bson:"destinations" json:"destinations"`
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator"`
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}
type BusinessMan struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	BusinessManDto `bson:",inline"`
}
