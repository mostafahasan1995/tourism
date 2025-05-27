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
    "image": [
        {
            "_id": "507f1f77bcf86cd799439014",
            "originalName": "profile1.jpg",
            "path": "/uploads/profile1.jpg",
            "service": "file-service",
            "expire": "2024-12-31T23:59:59Z",
            "variants": ["thumbnail", "medium", "large"]
        }
    ],
    "countries": ["USA", "France", "Italy", "Spain"],
    "contact": {
        "email": "john.smith@example.com",
        "phone": "+1234567890",
        "address": "123 Travel Street, New York, USA"
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

type AgentDto struct {
	Name        string            `bson:"name" json:"name"`
	Nationality string            `bson:"nationality" json:"nationality"`
	SpokenLangs []string          `bson:"languages" json:"languages"`
	Company     string            `bson:"company" json:"company"`
	CompanyLogo types.FileField   `bson:"companyLogo" json:"companyLogo"`
	Bio         string            `bson:"bio" json:"bio"`
	Image       []types.FileField `bson:"image" json:"image"`
	Countries   []string          `bson:"countries" json:"countries"`
	Contact     MemberContact     `bson:"contact" json:"contact"`
	Security    MemberSecurity    `bson:"security" json:"security"`
	Financial   AgentFinancial    `bson:"financial" json:"financial"`
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
