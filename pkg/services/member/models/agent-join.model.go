package models

import (
	"larsa-tourism-microservices/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentJoinDto struct {
	FullName    string          `bson:"fullName" json:"fullName"`
	Nationality string          `bson:"nationality" json:"nationality"`
	Bio         string          `bson:"bio" json:"bio"`
	Phone       string          `bson:"phone" json:"phone"`
	Email       string          `bson:"email" json:"email"`
	SpokenLangs []string        `bson:"spokenLangs" json:"spokenLangs"`
	CompanyName string          `bson:"companyName" json:"companyName"`
	CompanyLogo types.FileField `bson:"companyLogo" json:"companyLogo"`
	Countries   []string        `bson:"countries" json:"countries"`
}

type AgentJoin struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AgentJoinDto `bson:",inline"`
	Status       string `bson:"status,omitempty" json:"status,omitempty"` //pending - converted - rejected
}

type AgentJoinPagination struct {
	Agents     []AgentJoin      `bson:"agents" json:"agents"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
