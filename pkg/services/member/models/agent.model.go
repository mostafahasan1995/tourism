package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/member/enums"
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentDto struct {
	Name        transl.Localizable[string] `bson:"name" json:"name"`
	Nationality string                     `bson:"nationality" json:"nationality"`
	SpokenLangs []string                   `bson:"spokenLangs" json:"spokenLangs"`
	Company     transl.Localizable[string] `bson:"company" json:"company"`
	CompanyLogo types.FileField            `bson:"companyLogo" json:"companyLogo"`
	Bio         transl.Localizable[string] `bson:"bio" json:"bio"`
	Image       types.FileField            `bson:"image" json:"image"` // Changed from array to single object
	Countries   []string                   `bson:"countries" json:"countries"`
	Contact     AgentContact               `bson:"contacts" json:"contacts"` // Changed to new AgentContact structure
	Security    MemberSecurity             `bson:"security" json:"security"`
	Financial   AgentFinancial             `bson:"financial" json:"financial"`
}

type AgentFinancial struct {
	ProfitOfTourismProgram bool    `bson:"profitOfTourismProgram" json:"profitOfTourismProgram"`
	Ratio                  float64 `bson:"ratio" json:"ratio"`
}

type FinancialUnit struct {
	Active bool    `bson:"active" json:"active"`
	Cost   float64 `bson:"cost" json:"cost"`
}

type Agent struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // same as user id
	AgentId   string             `bson:"agentId,omitempty" json:"agentId,omitempty"`
	AgentDto  `bson:",inline"`
	Status    enums.AgentStatus  `bson:"status" json:"status"` // active, inactive
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

type UpdateStatusDto struct {
	Status string `json:"status" validate:"required,oneof=active inactive"`
}

func (u *UpdateStatusDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, u)
}
