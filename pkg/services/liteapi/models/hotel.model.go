package models

type AccessibilityAttributes struct {
	Attributes                                 any     `bson:"attributes" json:"attributes"`
	ShowerChair                                any     `bson:"showerChair" json:"showerChair"`
	EntranceType                               any     `bson:"entranceType" json:"entranceType"`
	PetFriendly                                string  `bson:"petFriendly" json:"petFriendly"`
	RampAngle                                  float64 `bson:"rampAngle" json:"rampAngle"`
	RampLength                                 float64 `bson:"rampLength" json:"rampLength"`
	EntranceDoorWidth                          float64 `bson:"entranceDoorWidth" json:"entranceDoorWidth"`
	RoomMaxGuestsNumber                        int     `bson:"roomMaxGuestsNumber" json:"roomMaxGuestsNumber"`
	DistanceFromTheElevatorToTheAccessibleRoom float64 `bson:"distanceFromTheElevatorToTheAccessibleRoom" json:"distanceFromTheElevatorToTheAccessibleRoom"`
}

type Hotel struct {
	Id                      string                  `bson:"id" json:"id"`
	PrimaryHotelId          any                     `bson:"primaryHotelId" json:"primaryHotelId"`
	Name                    string                  `bson:"name" json:"name"`
	HotelDescription        string                  `bson:"hotelDescription" json:"hotelDescription"`
	HotelTypeId             int                     `bson:"hotelTypeId" json:"hotelTypeId"`
	ChainId                 int                     `bson:"chainId" json:"chainId"`
	Chain                   string                  `bson:"chain" json:"chain"`
	Currency                string                  `bson:"currency" json:"currency"`
	Country                 string                  `bson:"country" json:"country"`
	City                    string                  `bson:"city" json:"city"`
	Latitude                float64                 `bson:"latitude" json:"latitude"`
	Longitude               float64                 `bson:"longitude" json:"longitude"`
	Address                 string                  `bson:"address" json:"address"`
	Zip                     string                  `bson:"zip" json:"zip"`
	MainPhoto               string                  `bson:"main_photo" json:"main_photo"`
	Thumbnail               string                  `bson:"thumbnail" json:"thumbnail"`
	Stars                   float64                 `bson:"stars" json:"stars"`
	Rating                  float64                 `bson:"rating" json:"rating"`
	ReviewCount             int                     `bson:"reviewCount" json:"reviewCount"`
	FacilityIds             []int                   `bson:"facilityIds" json:"facilityIds"`
	AccessibilityAttributes AccessibilityAttributes `bson:"accessibilityAttributes" json:"accessibilityAttributes"`
	DeletedAt               any                     `bson:"deletedAt" json:"deletedAt"`
}

//
