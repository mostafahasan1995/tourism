package models

/*
Example JSON structure for Agent:
{
    "_id": "507f1f77bcf86cd799439011",
    "agentId": "AGT123456",
    "name": "John Smith",
    "nationality": "United States",
    "languages": ["English", "Spanish", "French"],
    "company": "Global Travel Agency",
    "companyLogo": {
        "_id": "507f1f77bcf86cd799439013",
        "originalName": "company-logo.png",
        "path": "/uploads/company-logo.png",
        "service": "file-service",
        "expire": "2024-12-31T23:59:59Z",
        "variants": ["thumbnail", "medium", "large"]
    },
    "bio": "Experienced travel agent with 10 years in the industry",
    "image": {
        "_id": "507f1f77bcf86cd799439014",
        "originalName": "profile1.jpg",
        "path": "/uploads/profile1.jpg",
        "service": "file-service",
        "expire": "2024-12-31T23:59:59Z",
        "variants": ["thumbnail", "medium", "large"]
    },
    "countries": ["USA", "France", "Italy", "Spain"],
    "contacts": {
        "phone": {
            "pre": "+1",
            "content": "234567890"
        },
        "email": "john.smith@example.com",
        "web": "https://johntravelagent.com"
    },
    "security": {
        "password": "hashedPassword123",
        "role": "agent"
    },
    "financial": {
        "stays": {
            "active": true,
            "cost": 150.50
        },
        "tourismPrograms": {
            "active": true,
            "cost": 200.75
        },
        "vipCars": {
            "active": false,
            "cost": 0
        },
        "businessManTrip": {
            "active": true,
            "cost": 300.25
        },
        "delegation": {
            "active": true,
            "cost": 250.00
        }
    },
    "ratingObjects": [
        {
            "username": "customer1",
            "userId": "507f1f77bcf86cd799439015",
            "userImg": {
                "_id": "507f1f77bcf86cd799439016",
                "originalName": "user1.jpg",
                "path": "/uploads/user1.jpg",
                "service": "file-service",
                "expire": "2024-12-31T23:59:59Z",
                "variants": ["thumbnail", "medium", "large"]
            },
            "value": 4.5,
            "text": "Great agent, very helpful!",
            "status": "approved",
            "date": "2024-03-20T10:00:00Z",
            "replies": [
                {
                    "text": "Thank you for your feedback!",
                    "date": "2024-03-20T11:00:00Z"
                }
            ]
        }
    ],
    "ratings": 4.5,
    "status": "active",
    "isJoinReq": false,
    "joinStatus": "converted",
    "trash": false,
    "createdAt": "2024-03-20T10:00:00Z",
    "createdBy": "507f1f77bcf86cd799439012",
    "updatedAt": "2024-03-20T10:00:00Z",
    "updatedBy": "507f1f77bcf86cd799439012"
}
*/

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rating and review structures for agents (similar to hotels)
type AgentRatingObject struct {
	Username string          `bson:"username" json:"username"`
	UserId   string          `bson:"userId" json:"userId"`
	UserImg  types.FileField `bson:"userImg" json:"userImg"`
	Value    float64         `bson:"value" json:"value"`
	Text     string          `bson:"text" json:"text"`
	Status   string          `bson:"status" json:"status"`
	Date     *time.Time      `bson:"date" json:"date"`
	Replies  []AgentReply    `bson:"replies" json:"replies"`
}

type AgentReply struct {
	Text string     `bson:"text" json:"text"`
	Date *time.Time `bson:"date" json:"date"`
}

type AgentDto struct {
	Name          string              `bson:"name" json:"name"`
	Nationality   string              `bson:"nationality" json:"nationality"`
	SpokenLangs   []string            `bson:"Spokenlanguages" json:"Spokenlanguages"`
	Company       string              `bson:"company" json:"company"`
	CompanyLogo   types.FileField     `bson:"companyLogo" json:"companyLogo"`
	Bio           string              `bson:"bio" json:"bio"`
	Image         types.FileField     `bson:"image" json:"image"` // Changed from array to single object
	Countries     []string            `bson:"countries" json:"countries"`
	Contact       AgentContact        `bson:"contacts" json:"contacts"` // Changed to new AgentContact structure
	Security      MemberSecurity      `bson:"security" json:"security"`
	Financial     AgentFinancial      `bson:"financial" json:"financial"`
	RatingObjects []AgentRatingObject `bson:"ratingObjects" json:"ratingObjects"` // Added rating objects
	Ratings       float64             `bson:"ratings" json:"ratings"`             // Added average rating
}

// CalculateAverageRating calculates the average rating from RatingObjects
func (a *AgentDto) CalculateAverageRating() {
	if len(a.RatingObjects) == 0 {
		return
	}
	var total float64
	for _, rating := range a.RatingObjects {
		total += rating.Value
	}
	a.Ratings = total / float64(len(a.RatingObjects))
}

type AgentFinancial struct {
	Stays           FinancialUnit `bson:"stays" json:"stays"`
	TourismPrograms FinancialUnit `bson:"tourismPrograms" json:"tourismPrograms"`
	VipCars         FinancialUnit `bson:"vipCars" json:"vipCars"`
	BusinessManTrip FinancialUnit `bson:"businessManTrip" json:"businessManTrip"`
	Delegation      FinancialUnit `bson:"delegation" json:"delegation"`
}

type FinancialUnit struct {
	Active bool    `bson:"active" json:"active"`
	Cost   float64 `bson:"cost" json:"cost"`
}

type Agent struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // same as user id
	AgentId    string             `bson:"agentId,omitempty" json:"agentId,omitempty"`
	AgentDto   `bson:",inline"`
	Status     string             `bson:"status" json:"status"` // active, inactive
	IsJoinReq  bool               `bson:"isJoinReq" json:"isJoinReq"`
	JoinStatus string             `bson:"joinStatus" json:"joinStatus"` //converted, pending, rejected
	Trash      bool               `bson:"trash" json:"trash"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type AgentWithPagination struct {
	Agents     []Agent          `bson:"agents" json:"agents"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
