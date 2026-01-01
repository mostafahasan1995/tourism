package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FinancialSettingsDto struct {
	ProfitRatio float64     `bson:"profitRatio" json:"profitRatio"`
	BankAccount BankAccount `bson:"bankAccount" json:"bankAccount"`
}

type BankAccount struct {
	BankName          string   `bson:"bankName" json:"bankName"`
	AccountNumber     string   `bson:"accountNumber" json:"accountNumber"`
	AccountHolderName string   `bson:"accountHolderName" json:"accountHolderName"`
	IBAN              string   `bson:"iban" json:"iban"`
	SwiftCode         string   `bson:"swiftCode" json:"swiftCode"`
	Currencies        []string `bson:"currencies" json:"currencies"`
}

type FinancialSettings struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string
	FinancialSettingsDto `bson:",inline"`
}
