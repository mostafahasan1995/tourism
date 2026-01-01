package models

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

type HotelReviewList struct {
	Data              []HotelReview           `bson:"data" json:"data"`
	Total             int                     `bson:"total" json:"total"`
	SentimentAnalysis ReviewSentimentAnalysis `bson:"sentimentAnalysis,omitempty" json:"sentimentAnalysis,omitempty"`
}

type HotelReview struct {
	Id           string    `bson:"id" json:"-"`
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

func (hr *HotelReview) GenerateReviewID() {
	// Format date consistently to ensure same output for same date
	dateStr := hr.Date.UTC().Format("2006-01-02")

	// Create a consistent string combining all elements
	// Using | as separator since it's unlikely to appear in names
	baseString := fmt.Sprintf("%.2f|%s|%s|%s|%s",
		hr.AverageScore,
		strings.ToLower(strings.TrimSpace(hr.Name)),
		dateStr,
		strings.ToLower(strings.TrimSpace(hr.Type)),
		strings.ToLower(strings.TrimSpace(hr.Language)),
	)

	// Create a hash of the string
	h := fnv.New32a()
	h.Write([]byte(baseString))
	hash := h.Sum32()

	// Convert to base36 for shorter, readable string
	// This will give us alphanumeric IDs
	id := fmt.Sprintf("r%s", strings.ToUpper(strconv.FormatUint(uint64(hash), 36)))
	hr.Id = id
}

type ReviewSentimentAnalysis struct {
	HotelId           string `bson:"hotelId" json:"-"`
	SentimentAnalysis `bson:",inline"`
}
