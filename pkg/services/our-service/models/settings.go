package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	ProfitRatio float64            `bson:"profitRatio" json:"profitRatio"`
}
