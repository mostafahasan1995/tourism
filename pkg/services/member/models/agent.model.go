package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentDto struct {
	Name        string          `bson:"name" json:"name"`
	Nationality string          `bson:"nationality" json:"nationality"`
	SpokenLangs []string        `bson:"spokenLangs" json:"spokenLangs"`
	Company     string          `bson:"company" json:"company"`
	CompanyLogo types.FileField `bson:"companyLogo" json:"companyLogo"`
	Bio         string          `bson:"bio" json:"bio"`
	Image       types.FileField `bson:"image" json:"image"` // Changed from array to single object
	Countries   []string        `bson:"countries" json:"countries"`
	Contact     AgentContact    `bson:"contacts" json:"contacts"` // Changed to new AgentContact structure
	Security    MemberSecurity  `bson:"security" json:"security"`
	Financial   AgentFinancial  `bson:"financial" json:"financial"`
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
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // same as user id
	AgentId   string             `bson:"agentId,omitempty" json:"agentId,omitempty"`
	AgentDto  `bson:",inline"`
	Status    string             `bson:"status" json:"status"` // active, inactive
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type AgentWithPagination struct {
	Agents     []Agent          `bson:"agents" json:"agents"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
