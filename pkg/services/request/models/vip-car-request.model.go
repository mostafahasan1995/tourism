package models

import (
	// "larsa-tourism-microservices/pkg/types"
	// "larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VipCarRequestDto struct {
	Location              string   `bson:"location" json:"location"`
	Capacity              string   `bson:"capacity" json:"capacity"`
	DriverLanguagesSpoken []string `bson:"driverLanguagesSpoken" json:"driverLanguagesSpoken"`
	LuxuryFeatures        []string `bson:"luxuryFeatures" json:"luxuryFeatures"`

	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`

	StartTime string `bson:"startTime" json:"startTime"`
	EndTime   string `bson:"endTime" json:"endTime"`

	CarTypeId primitive.ObjectID`bson:"carTypeId" json:"carTypeId"`
}

type VipCarRequest struct {
	VipCarRequestDto `bson:",inline"`

	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Trash bool `bson:"trash" json:"trash"`

	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type VipCarRequestPagination struct {
	VipCarRequest []VipCarRequest `bson:"vipCarRequest" json:"vipCarRequest"`

	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
