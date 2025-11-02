package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// general program - to show in the website so the customer can book it then we will create a custom program for them
type GeneralProgram struct {
	Destinations   []GeneralProgramDestination `bson:"destinations" json:"destinations"`
	Includes       Includes                    `bson:"includes" json:"includes"`
	Activities     any                         `bson:"activities" json:"activities"`
	DailyItinerary []DailyItinerary            `bson:"dailyItinerary" json:"dailyItinerary"`
	Notes          GPNotes                     `bson:"notes" json:"notes"`
	Pricing        GPPricing                   `bson:"pricing" json:"pricing"`
}

type GPNotes struct {
	SutableFor     transl.Localizable[string] `bson:"sutableFor" json:"sutableFor"`
	ImportantNotes transl.Localizable[string] `bson:"importantNotes" json:"importantNotes"`
}

type GeneralProgramDestination struct {
	From        primitive.ObjectID         `bson:"from" json:"from"`
	Images      []types.FileField          `bson:"images" json:"images"`
	Duration    int                        `bson:"duration" json:"duration"`
	Description transl.Localizable[string] `bson:"description" json:"description"`
}

type Includes struct {
	Accommodation      []string `bson:"accommodation" json:"accommodation"`
	Transportation     []string `bson:"transportation" json:"transportation"`
	Meals              []string `bson:"meals" json:"meals"`
	IncludedDetails    string   `bson:"includedDetails" json:"includedDetails"`
	NotIncludedDetails string   `bson:"notIncludedDetails" json:"notIncludedDetails"`
}

type DailyItinerary struct {
	Title      transl.Localizable[string]   `bson:"title" json:"title"`
	Actions    []primitive.ObjectID         `bson:"actions" json:"actions"`
	NewActions []transl.Localizable[string] `bson:"newActions" json:"newActions"`
	Images     []types.FileField            `bson:"images" json:"images"`
}

type GPPricing struct {
	Person   PersonPrice   `bson:"person" json:"person"`
	Children ChildrenPrice `bson:"children" json:"children"`
}

type PersonPrice struct {
	Price         float64 `bson:"price" json:"price"`
	Per           string  `bson:"per" json:"per"`
	ShowInWebsite bool    `bson:"showInWebsite" json:"showInWebsite"`
}
type ChildrenPrice struct {
	Price         float64 `bson:"price" json:"price"`
	Per           string  `bson:"per" json:"per"`
	NumOfYears    string  `bson:"numOfYears" json:"numOfYears"`
	ShowInWebsite bool    `bson:"showInWebsite" json:"showInWebsite"`
}
