package models

import "time"

type HotelReviewList struct {
	Data              []HotelReview           `bson:"data" json:"data"`
	Total             int                     `bson:"total" json:"total"`
	SentimentAnalysis ReviewSentimentAnalysis `bson:"sentimentAnalysis,omitempty" json:"sentimentAnalysis,omitempty"`
}

type HotelReview struct {
	Seq          int       `bson:"seq" json:"-"`
	HotelId      string    `bson:"hotelId" json:"-"`
	AverageScore float64   `bson:"averageScore" json:"averageScore"`
	Country      string    `bson:"country" json:"country"`
	Type         string    `bson:"type" json:"type"`
	Name         string    `bson:"name" json:"name"`
	Date         time.Time `bson:"date" json:"date"`
	Headline     string    `bson:"headline" json:"headline"`
	Language     string    `bson:"language" json:"language"`
	Pros         string    `bson:"pros" json:"pros"`
	Cons         string    `bson:"cons" json:"cons"`
	Source       string    `bson:"source" json:"source"`
	ExpiresAt    time.Time `bson:"expiresAt" json:"-"`
}

type ReviewSentimentAnalysis struct {
	HotelId           string `bson:"hotelId" json:"-"`
	SentimentAnalysis `bson:",inline"`
}
