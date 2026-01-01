package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ManualStatistics represents manually entered statistics data
type ManualStatistics struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProgramsCount      int64              `bson:"programsCount" json:"programsCount"`
	TripsCount         int64              `bson:"tripsCount" json:"tripsCount"`
	CountriesCount     int64              `bson:"countriesCount" json:"countriesCount"`
	HotelsCount        int64              `bson:"hotelsCount" json:"hotelsCount"`
	HappyTravelerCount int64              `bson:"happyTravelerCount" json:"happyTravelerCount"`
	AutoCalculate      bool               `bson:"autoCalculate" json:"autoCalculate"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// ManualStatisticsRequest represents the request payload for updating manual statistics
type ManualStatisticsRequest struct {
	ProgramsCount      *int64 `json:"programsCount,omitempty"`
	TripsCount         *int64 `json:"tripsCount,omitempty"`
	CountriesCount     *int64 `json:"countriesCount,omitempty"`
	HotelsCount        *int64 `json:"hotelsCount,omitempty"`
	HappyTravelerCount *int64 `json:"happyTravelerCount,omitempty"`
	AutoCalculate      *bool  `json:"autoCalculate,omitempty"`
}

// AutoCalculateRequest represents the request to toggle auto-calculate mode
type AutoCalculateRequest struct {
	AutoCalculate bool `json:"autoCalculate"`
}
