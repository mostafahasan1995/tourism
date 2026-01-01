package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Place struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PlaceId  string             `bson:"placeId" json:"placeId"`
	Language string             `bson:"language" json:"language"`
	Data     SearchPlaceData    `bson:"data" json:"data"`
}

type PlaceResponse struct {
	Data SearchPlaceData `bson:"data" json:"data"`
}

type SearchPlaceData struct {
	AddressComponents []AddressComponent `bson:"addressComponents" json:"addressComponents"`
	Location          Coordinates        `bson:"location" json:"location"`
	Viewport          Viewport           `bson:"viewport" json:"viewport"`
}

type AddressComponent struct {
	LanguageCode string   `bson:"languageCode" json:"languageCode"`
	LongText     string   `bson:"longText" json:"longText"`
	ShortText    string   `bson:"shortText" json:"shortText"`
	Types        []string `bson:"types" json:"types"`
}

type Coordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

type Viewport struct {
	High ViewportBounds `bson:"high" json:"high"`
	Low  ViewportBounds `bson:"low" json:"low"`
}

type ViewportBounds struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}
