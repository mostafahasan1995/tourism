package models

import (
	"larsa-tourism-microservices/pkg/helpers"
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
	ServiceType string             `bson:"serviceType" json:"serviceType"`         // e.g. delegation - custom-plan - business-man - vip-car - flight-request - partner-request
	TravelReqId primitive.ObjectID `bson:"travelReqId,omitempty" json:"travelReqId,omitempty"`
	CustomerId  primitive.ObjectID `bson:"customerId,omitempty" json:"customerId,omitempty"`
	Status      string             `bson:"status" json:"status"`
	Package     primitive.ObjectID `bson:"package" json:"package" validate:"required"`
	ProgramType string             `bson:"programType" json:"programType" validate:"required"` // general - custom
	//
	Source      string    `bson:"source" json:"source"`
	Company     string    `bson:"company" json:"company"`         // todo: maybe we need id here
	Coordinator string    `bson:"coordinator" json:"coordinator"` // todo: maybe we need id here
	Purpose     string    `bson:"purpose" json:"purpose"`
	StartDate   time.Time `bson:"startDate" json:"startDate"`
	EndDate     time.Time `bson:"endDate" json:"endDate"`
	GroupSize   string    `bson:"groupSize" json:"groupSize"` //see group size values above
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

/*
Example JSON structure for ProgramDto (General Program):
{
    "title": "Luxury Paris Experience",
    "serviceType": "custom-plan",
    "travelReqId": "507f1f77bcf86cd799439011",
    "customerId": "507f1f77bcf86cd799439012",
    "status": "active",
    "package": "507f1f77bcf86cd799439013",
    "programType": "general",
    "source": "website",
    "company": "Luxury Travel Co.",
    "coordinator": "John Smith",
    "purpose": "Luxury Vacation",
    "startDate": "2024-06-01T00:00:00Z",
    "endDate": "2024-06-07T00:00:00Z",
    "groupSize": "Couple",
    "generalType": {
        "distinations": ["507f1f77bcf86cd799439014"],
        "includes": {
            "accommodation": ["5-star hotel", "Luxury suite"],
            "transportation": ["Private car", "Airport transfer"],
            "meals": ["Breakfast", "Welcome dinner"]
        },
        "activities": ["507f1f77bcf86cd799439015"],
        "dailyItinerary": [
            {
                "title": "Day 1 - Arrival in Paris",
                "actions": ["507f1f77bcf86cd799439016"],
                "images": [
                    {
                        "_id": "507f1f77bcf86cd799439017",
                        "originalName": "welcome-dinner.jpg",
                        "path": "/uploads/welcome-dinner.jpg",
                        "service": "file-service",
                        "expire": "2024-12-31T23:59:59Z",
                        "variants": ["thumbnail", "medium", "large"]
                    }
                ]
            }
        ],
        "pricing": {
            "person": {
                "price": 5000,
                "per": "person",
                "showInWebsite": true
            },
            "children": {
                "price": 2500,
                "per": "child",
                "numOfYears": "2-12",
                "showInWebsite": true
            }
        }
    }
}

Example JSON structure for ProgramDto (Custom Program):
{
    "title": "Corporate Team Building Retreat",
    "serviceType": "custom-plan",
    "travelReqId": "507f1f77bcf86cd799439021",
    "customerId": "507f1f77bcf86cd799439022",
    "status": "active",
    "package": "507f1f77bcf86cd799439023",
    "programType": "custom",
    "source": "direct",
    "company": "Tech Solutions Inc.",
    "coordinator": "Sarah Johnson",
    "purpose": "Team Building and Strategy Planning",
    "startDate": "2024-08-01T00:00:00Z",
    "endDate": "2024-08-05T00:00:00Z",
    "groupSize": "Small Group",
    "customType": {
        "customPlan": {
            "tripType": "luxury retreat",
            "nationality": "American",
            "tripDuration": 5,
            "tripCoordinator": "507f1f77bcf86cd799439026"
        },
        "destinations": [
            {
                "destination": "507f1f77bcf86cd799439024",
                "tripDetails": {
                    "startDate": "2024-08-01T10:00:00Z",
                    "endDate": "2024-08-05T14:00:00Z",
                    "days": 5
                },
                "accommodation": [
                    {
                        "city": "Bali",
                        "preferences": "Resort with meeting facilities",
                        "specialRequests": ["Ocean view", "Meeting rooms"],
                        "numOfRooms": 5,
                        "bedType": "King",
                        "numOfBeds": 1,
                        "numOfBathrooms": 1,
                        "privateMeetingRoom": true
                    }
                ],
                "flightTickets": {
                    "arrangeByUs": true,
                    "tripType": "round-trip",
                    "travelClass": "business",
                    "numOfPassengers": 8,
                    "departureDate": "2024-08-01T08:00:00Z",
                    "returnDate": "2024-08-05T16:00:00Z",
                    "flexibleTravelDates": "No",
                    "preferredDepartureTime": "morning",
                    "layoverPreferences": "short transit time",
                    "preferredAirlines": "Singapore Airlines",
                    "extraBaggage": true,
                    "specialMeals": true,
                    "travelWithPet": false,
                    "anySpecialReq": "Group seating preferred",
                    "adults": 8,
                    "children": 0,
                    "infant": 0
                },
                "transportation": {
                    "trans": "Private van",
                    "driverLangs": ["English", "Indonesian"]
                },
                "activities": ["Team building", "Strategy workshop", "Cultural tour"],
                "agenda": {
                    "needAgenda": true,
                    "items": [
                        "Morning team building",
                        "Afternoon strategy planning",
                        "Evening cultural activities"
                    ],
                    "files": [
                        {
                            "_id": "507f1f77bcf86cd799439025",
                            "originalName": "team-building-schedule.pdf",
                            "path": "/uploads/team-building-schedule.pdf",
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
                    "tourAfterMeeting": true,
                    "photography": true,
                    "airportMeetAndGreet": true,
                    "simCardAndInternet": true
                }
            }
        ]
    }
}
*/
