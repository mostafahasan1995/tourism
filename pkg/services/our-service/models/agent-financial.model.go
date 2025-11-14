package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentWithdrawal struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Amount float64            `bson:"amount" json:"amount"`
	Method string             `bson:"method" json:"method"`
	Note   string             `bson:"note" json:"note"`
	Date   time.Time          `bson:"date" json:"date"`
}

type AgentFinancialAccount struct {
	Id             primitive.ObjectID         `bson:"_id,omitempty" json:"_id,omitempty"`
	AgentId        primitive.ObjectID         `bson:"agentId" json:"agentId"`
	AgentName      transl.Localizable[string] `bson:"agentName" json:"agentName"`
	TotalProfit    float64                    `bson:"totalProfit" json:"totalProfit"`
	TotalWithdrawn float64                    `bson:"totalWithdrawn" json:"totalWithdrawn"`
	Balance        float64                    `bson:"balance" json:"balance"`
	Withdrawals    []AgentWithdrawal          `bson:"withdrawals" json:"withdrawals"`
	CreatedAt      time.Time                  `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      time.Time                  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type AgentWithdrawRequest struct {
	Amount float64 `json:"amount"`
	Method string  `json:"method"`
	Note   string  `json:"note"`
}
