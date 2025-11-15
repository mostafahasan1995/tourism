package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/transl"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentWithdrawal struct {
	Id     primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
	Amount float64                `bson:"amount" json:"amount"`
	Method string                 `bson:"method" json:"method"`
	Note   string                 `bson:"note" json:"note"`
	Status enums.WithdrawalStatus `bson:"status" json:"status"`
	Date   time.Time              `bson:"date" json:"date"`
}

type AgentFinancialAccount struct {
	Id              primitive.ObjectID         `bson:"_id,omitempty" json:"_id,omitempty"`
	AgentId         primitive.ObjectID         `bson:"agentId" json:"agentId"`
	AgentName       transl.Localizable[string] `bson:"agentName" json:"agentName"`
	PaymentMethod   string                     `bson:"paymentMethod" json:"paymentMethod"`     // e.g., "bank-transfer", "cash", "paypal", etc.
	BankAccountInfo string                     `bson:"bankAccountInfo" json:"bankAccountInfo"` // Bank address, account number, IBAN, etc.
	TotalProfit     float64                    `bson:"totalProfit" json:"totalProfit"`
	TotalWithdrawn  float64                    `bson:"totalWithdrawn" json:"totalWithdrawn"`
	Balance         float64                    `bson:"balance" json:"balance"`
	Withdrawals     []AgentWithdrawal          `bson:"withdrawals" json:"withdrawals"`
	CreatedAt       time.Time                  `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt       time.Time                  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type AgentFinancialAccountDto struct {
	PaymentMethod   string `json:"paymentMethod" validate:"required"`
	BankAccountInfo string `json:"bankAccountInfo"`
}

func (a *AgentFinancialAccountDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, a)
}

type AgentWithdrawRequest struct {
	Amount float64 `json:"amount"`
	Method string  `json:"method"`
	Note   string  `json:"note"`
}
