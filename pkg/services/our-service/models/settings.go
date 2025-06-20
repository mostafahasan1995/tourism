package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SettingsDto struct {
	ProfitRatio float64 `bson:"profitRatio" json:"profitRatio"`
}

type Settings struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string
	SettingsDto `bson:",inline"`
}
