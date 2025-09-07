package models

import "time"

type HotelReview struct {
	HotelId       string    `bson:"hotelId" json:"hotelId"`
	AverageScore  float64   `bson:"averageScore" json:"averageScore"`
	Country       string    `bson:"country" json:"country"`
	Type          string    `bson:"type" json:"type"`
	Name          string    `bson:"name" json:"name"`
	Date          time.Time `bson:"date" json:"date"`
	Headline      string    `bson:"headline" json:"headline"`
	Language      string    `bson:"language" json:"language"`
	Pros          string    `bson:"pros" json:"pros"`
	Cons          string    `bson:"cons" json:"cons"`
	Source        string    `bson:"source" json:"source"`
	LastUpdatedAt time.Time `bson:"lastUpdatedAt" json:"lastUpdatedAt"`
}

// Example JSON:
// {
// 	"averageScore": 8,
// 	"country": "de",
// 	"type": "family_with_children",
// 	"name": "Markus",
// 	"date": "2025-08-17T00:00:00Z",
// 	"headline": "Die Lage direkt am Time Square ist unschlagbar",
// 	"language": "de",
// 	"pros": "Die Zentrale Lage",
// 	"cons": "Das Frühstück war nur ein Lunch Paket",
// 	"source": "Nuitee"
// }
