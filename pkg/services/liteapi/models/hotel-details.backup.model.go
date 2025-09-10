package models

import (
	"time"
)

// HotelDetails represents the complete hotel details response
type HotelDetail struct {
	Id                        string               `json:"id" bson:"id"`
	Name                      string               `json:"name" bson:"name"`
	HotelDescription          string               `json:"hotelDescription" bson:"hotelDescription"`
	HotelImportantInformation string               `json:"hotelImportantInformation" bson:"hotelImportantInformation"`
	CheckinCheckoutTimes      CheckinCheckoutTimes `json:"checkinCheckoutTimes" bson:"checkinCheckoutTimes"`
	HotelImages               []HotelImage         `json:"hotelImages" bson:"hotelImages"`
	MainPhoto                 string               `json:"main_photo" bson:"main_photo"`
	Thumbnail                 string               `json:"thumbnail" bson:"thumbnail"`
	Country                   string               `json:"country" bson:"country"`
	City                      string               `json:"city" bson:"city"`
	StarRating                float64              `json:"starRating" bson:"starRating"`
	Location                  Location             `json:"location" bson:"location"`
	Address                   string               `json:"address" bson:"address"`
	HotelFacilities           []string             `json:"hotelFacilities" bson:"hotelFacilities"`
	Zip                       string               `json:"zip" bson:"zip"`
	Chain                     string               `json:"chain" bson:"chain"`
	Facilities                []Facility           `json:"facilities" bson:"facilities"`
	Rooms                     []Room               `json:"rooms" bson:"rooms"`
	Phone                     string               `json:"phone" bson:"phone"`
	Fax                       string               `json:"fax" bson:"fax"`
	Email                     string               `json:"email" bson:"email"`
	HotelType                 string               `json:"hotelType" bson:"hotelType"`
	HotelTypeID               int                  `json:"hotelTypeId" bson:"hotelTypeId"`
	AirportCode               string               `json:"airportCode" bson:"airportCode"`
	Rating                    float64              `json:"rating" bson:"rating"`
	ReviewCount               int                  `json:"reviewCount" bson:"reviewCount"`
	Parking                   string               `json:"parking" bson:"parking"`
	GroupRoomMin              int                  `json:"groupRoomMin" bson:"groupRoomMin"`
	ChildAllowed              bool                 `json:"childAllowed" bson:"childAllowed"`
	PetsAllowed               bool                 `json:"petsAllowed" bson:"petsAllowed"`
	Policies                  []Policy             `json:"policies" bson:"policies"`
	SentimentAnalysis         SentimentAnalysis    `json:"sentiment_analysis" bson:"sentiment_analysis"`
	SentimentUpdatedAt        time.Time            `json:"sentiment_updated_at" bson:"sentiment_updated_at"`
	DeletedAt                 *time.Time           `json:"deletedAt" bson:"deletedAt"`
	ExpiresAt                 time.Time            `bson:"expiresAt" json:"-"`
}

// CheckinCheckoutTimes represents check-in and check-out times
type CheckinCheckoutTimes struct {
	Checkout string `json:"checkout" bson:"checkout"`
	Checkin  string `json:"checkin" bson:"checkin"`
}

// HotelImage represents a hotel image
type HotelImage struct {
	URL          string `json:"url" bson:"url"`
	URLHd        string `json:"urlHd" bson:"urlHd"`
	Caption      string `json:"caption" bson:"caption"`
	Order        int    `json:"order" bson:"order"`
	DefaultImage bool   `json:"defaultImage" bson:"defaultImage"`
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// Facility represents a hotel facility
type Facility struct {
	FacilityID int    `json:"facilityId" bson:"facilityId"`
	Name       string `json:"name" bson:"name"`
}

// Room represents a hotel room
type Room struct {
	Id             int           `json:"id" bson:"id"`
	RoomName       string        `json:"roomName" bson:"roomName"`
	Description    string        `json:"description" bson:"description"`
	RoomSizeSquare float64       `json:"roomSizeSquare" bson:"roomSizeSquare"`
	RoomSizeUnit   string        `json:"roomSizeUnit" bson:"roomSizeUnit"`
	HotelID        string        `json:"hotelId" bson:"hotelId"`
	MaxAdults      int           `json:"maxAdults" bson:"maxAdults"`
	MaxChildren    int           `json:"maxChildren" bson:"maxChildren"`
	MaxOccupancy   int           `json:"maxOccupancy" bson:"maxOccupancy"`
	BedTypes       []BedType     `json:"bedTypes" bson:"bedTypes"`
	RoomAmenities  []RoomAmenity `json:"roomAmenities" bson:"roomAmenities"`
	Photos         []RoomPhoto   `json:"photos" bson:"photos"`
	Views          []string      `json:"views" bson:"views"`
	BedRelation    string        `json:"bedRelation" bson:"bedRelation"`
}

// BedType represents bed type information
type BedType struct {
	Quantity int    `json:"quantity" bson:"quantity"`
	BedType  string `json:"bedType" bson:"bedType"`
	BedSize  string `json:"bedSize" bson:"bedSize"`
	ID       int    `json:"id" bson:"id"`
}

// RoomAmenity represents room amenity information
type RoomAmenity struct {
	AmenitiesID int    `json:"amenitiesId" bson:"amenitiesId"`
	Name        string `json:"name" bson:"name"`
	Sort        int    `json:"sort" bson:"sort"`
}

// RoomPhoto represents room photo information
type RoomPhoto struct {
	URL              string  `json:"url" bson:"url"`
	ImageDescription string  `json:"imageDescription" bson:"imageDescription"`
	ImageClass1      string  `json:"imageClass1" bson:"imageClass1"`
	ImageClass2      string  `json:"imageClass2" bson:"imageClass2"`
	FailoverPhoto    string  `json:"failoverPhoto" bson:"failoverPhoto"`
	MainPhoto        bool    `json:"mainPhoto" bson:"mainPhoto"`
	Score            float64 `json:"score" bson:"score"`
	ClassID          int     `json:"classId" bson:"classId"`
	ClassOrder       int     `json:"classOrder" bson:"classOrder"`
	HdURL            string  `json:"hd_url" bson:"hd_url"`
}

// Policy represents hotel policy information
type Policy struct {
	Id           int    `json:"id" bson:"id"`
	PolicyType   string `json:"policy_type" bson:"policy_type"`
	Name         string `json:"name" bson:"name"`
	Description  string `json:"description" bson:"description"`
	ChildAllowed string `json:"child_allowed" bson:"child_allowed"`
	PetsAllowed  string `json:"pets_allowed" bson:"pets_allowed"`
	Parking      string `json:"parking" bson:"parking"`
}

// SentimentAnalysis represents sentiment analysis data
type SentimentAnalysis struct {
	Cons       []string            `json:"cons" bson:"cons"`
	Pros       []string            `json:"pros" bson:"pros"`
	Categories []SentimentCategory `json:"categories" bson:"categories"`
}

// SentimentCategory represents a sentiment category
type SentimentCategory struct {
	Name        string  `json:"name" bson:"name"`
	Rating      float64 `json:"rating" bson:"rating"`
	Description string  `json:"description" bson:"description"`
}
