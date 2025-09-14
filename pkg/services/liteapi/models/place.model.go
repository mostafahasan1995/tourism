package models

type Place struct {
	PlaceId          string   `bson:"placeId" json:"placeId"`
	DisplayName      string   `bson:"displayName" json:"displayName"`
	FormattedAddress string   `bson:"formattedAddress" json:"formattedAddress"`
	Types            []string `bson:"types" json:"types"`
	Language         string   `bson:"language" json:"language"`
}
