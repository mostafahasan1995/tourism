package models

import (
	// "larsa-tourism-microservices/pkg/types"
	// "larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VipCarRequestDto struct {
	Destinations          []primitive.ObjectID `bson:"destinations" json:"destinations"`
	Capacity              string               `bson:"capacity" json:"capacity"`
	DriverLanguagesSpoken []string             `bson:"driverLanguagesSpoken" json:"driverLanguagesSpoken"`
	LuxuryFeatures        []string             `bson:"luxuryFeatures" json:"luxuryFeatures"`

	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	StartTime string    `bson:"startTime" json:"startTime"`
	EndTime   string    `bson:"endTime" json:"endTime"`

	CarTypeId primitive.ObjectID `bson:"carTypeId" json:"carTypeId"`
}

type VipCarRequest struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	VipCarRequestDto `bson:",inline"`
}
