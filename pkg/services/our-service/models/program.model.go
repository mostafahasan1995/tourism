package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Travel Type values:
// Relaxation Trip (رحلة استجمام)
// Adventure (مغامرة)
// Family Trip (رحلة عائلية)
// Romantic Trip (Honeymoon) (رحلة رومانسية – شهر عسل)
// Cultural Trip (رحلة ثقافية)
// Business Trip (رحلة عمل)
// Shopping Trip (رحلة تسوق)
// Wellness or Medical Tourism (رحلة صحية أو استشفائية)

// ✅ عدد الأفراد (Group Size):
// Solo Traveler (فردي)s
// Couple (زوجان)
// Family (عائلة)
// Small Group (مجموعة صغيرة، عادة 4–8 أشخاص)
// Large Group (مجموعة كبيرة، عادة أكثر من 8 أشخاص)

// Example JSON for a custom business-man program:
// {
//   "title": "Business Trip to Dubai",
//   "serviceType": "business-man",
//   "package": "507f1f77bcf86cd799439011",
//   "status": "pending",
//   "programType": "custom",
//   "source": "website",
//   "company": "Tech Corp",
//   "coordinator": "John Smith",
//   "purpose": "Business Meeting",
//   "startDate": "2024-03-20T00:00:00Z",
//   "endDate": "2024-03-25T00:00:00Z",
//   "groupSize": "solo",
//   "customType": {
//     "businessMan": {
//       "purpose": "meeting",
//       "clientIsCoordinator": false,
//       "coordinatorName": "Sarah Johnson",
//       "coordinatorPhone": "+1987654321",
//       "coordinatorEmail": "sarah.j@company.com"
//     },
//     "destinations": [
//       {
//         "destinationFrom": "507f1f77bcf86cd799439011",
//         "destinationTo": "507f1f77bcf86cd799439012",
//         "tripDetails": {
//           "startDate": "2024-03-20T10:00:00Z",
//           "endDate": "2024-03-25T14:00:00Z",
//           "days": 5
//         },
//         "accommodation": [
//           {
//             "city": "Dubai",
//             "preferences": "Business hotel in downtown",
//             "specialRequests": ["High floor", "Quiet room"],
//             "numOfRooms": 1,
//             "bedType": "King",
//             "numOfBeds": 1,
//             "numOfBathrooms": 1,
//             "privateMeetingRoom": true,
//             "pricePerNight": 250.00,
//             "totalStayCost": 1250.00
//           }
//         ],
//         "flightTickets": {
//           "arrangeByUs": true,
//           "tripType": "round-trip",
//           "travelClass": "business",
//           "departureDate": "2024-03-20T08:00:00Z",
//           "returnDate": "2024-03-25T16:00:00Z",
//           "flexibleTravelDates": "No",
//           "preferredDepartureTime": "morning",
//           "layoverPreferences": "short transit time",
//           "preferredAirlines": "Emirates",
//           "extraBaggage": true,
//           "specialMeals": true,
//           "travelWithPet": false,
//           "adults": 1,
//           "children": 0,
//           "infant": 0,
//           "totalCost": 1500.00
//         },
//         "transportation": {
//           "transType": "private",
//           "capacity": "4 passengers",
//           "driverLanguagesSpoken": ["English", "Arabic"],
//           "luxuryFeatures": ["Leather seats", "GPS", "Bluetooth", "Climate control"],
//           "startDate": "2024-03-20T10:00:00Z",
//           "endDate": "2024-03-25T14:00:00Z",
//           "startTime": "10:00",
//           "endTime": "14:00",
//           "carTypeId": "507f1f77bcf86cd799439013",
//           "totalCost": 500.00
//         },
//         "activities": {
//           "activities": ["City tour", "Business networking"],
//           "totalCost": 300.00
//         },
//         "agenda": {
//           "needAgenda": true,
//           "items": ["Morning meetings", "Afternoon site visits"],
//           "files": []
//         },
//         "services": {
//           "tourGuide": {
//             "active": true,
//             "cost": 200.00
//           },
//           "translator": {
//             "active": true,
//             "cost": 150.00
//           },
//           "airportPickup": {
//             "active": true,
//             "cost": 100.00
//           },
//           "tourAfterMeeting": {
//             "active": true,
//             "cost": 250.00
//           },
//           "photography": {
//             "active": false,
//             "cost": 0.00
//           },
//           "airportMeetAndGreet": {
//             "active": true,
//             "cost": 75.00
//           },
//           "simCardAndInternet": {
//             "active": true,
//             "cost": 50.00
//           }
//         }
//       }
//     ]
//   }
// }

type Program struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProgramDto `bson:",inline"`
	Trash      bool               `bson:"trash" json:"trash"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type ProgramDto struct {
	Title       string             `bson:"title" json:"title" validate:"required"` // program title
	ServiceType enums.ServiceType  `bson:"serviceType" json:"serviceType"`         // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	TravelReqId primitive.ObjectID `bson:"travelReqId" json:"travelReqId"`
	CustomerId  primitive.ObjectID `bson:"customerId" json:"customerId"`
	Status      string             `bson:"status" json:"status" validate:"required,oneof=pending active unactive"`
	Package     primitive.ObjectID `bson:"package" json:"package" validate:"required"`
	ProgramType string             `bson:"programType" json:"programType" validate:"required,oneof=general custom"` // general - custom
	//
	Source      string          `bson:"source" json:"source"`
	Company     string          `bson:"company" json:"company"`         // todo: maybe we need id here
	Coordinator string          `bson:"coordinator" json:"coordinator"` // todo: maybe we need id here
	Purpose     string          `bson:"purpose" json:"purpose"`
	StartDate   time.Time       `bson:"startDate" json:"startDate"`
	EndDate     time.Time       `bson:"endDate" json:"endDate"`
	GroupSize   enums.GroupSize `bson:"groupSize" json:"groupSize"` //see group size values above
	//
	GeneralType *GeneralProgram `bson:"generalType,omitempty" json:"generalType,omitempty" validate:"required_if=ProgramType general"`
	CustomType  *CustomProgram  `bson:"customType,omitempty" json:"customType,omitempty" validate:"required_if=ProgramType custom"`
}

func (p *ProgramDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, p)
}

type ProgramPagination struct {
	Programs   []Program        `bson:"programs" json:"programs"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
