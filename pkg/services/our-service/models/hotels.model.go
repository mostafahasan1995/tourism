package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelsDto struct {
	Name                          string                        `bson:"name" json:"name"`
	HotelType                     string                        `bson:"hotelType" json:"hotelType"`
	Location                      string                        `bson:"location" json:"location"`
	CheckInAndCheckOut            CheckInAndCheckOut            `bson:"checkInAndCheckOut" json:"checkInAndCheckOut"`
	Price                         int                           `bson:"price" json:"price"`
	IsDisplayInPerfectStay        bool                          `bson:"isDisplayInPerfectStay" json:"isDisplayInPerfectStay"`
	Ratings                       int                           `bson:"ratings" json:"ratings"`
	Image                         types.FileField               `bson:"image" json:"image"`
	RoomAmenities                 []string                      `bson:"roomAmenities" json:"roomAmenities"`
	DistanceFromCityCenter        int                           `bson:"distanceFromCityCenter" json:"distanceFromCityCenter"`
	NearbyAttractions             []string                      `bson:"nearbyAttractions" json:"nearbyAttractions"`
	ImagesGallery                 []types.FileField             `bson:"imagesGallery" json:"imagesGallery"`
	OverviewPage                  OverviewPage                  `bson:"overviewPage" json:"overviewPage"`
	RoomsAndSuitesPage            RoomsAndSuitesPage            `bson:"roomsAndSuitesPage" json:"roomsAndSuitesPage"`
	AmenitiesAndFacilitiesPage    AmenitiesAndFacilitiesPage    `bson:"amenitiesAndFacilitiesPage" json:"amenitiesAndFacilitiesPage"`
	LocationNearbyAttractionsPage LocationNearbyAttractionsPage `bson:"locationNearbyAttractionsPage" json:"locationNearbyAttractionsPage"`
	ReviewsAndRatingsPage         ReviewsAndRatingsPage         `bson:"reviewsAndRatingsPage" json:"reviewsAndRatingsPage"`
	BookingAndPoliciesPage        BookingAndPoliciesPage        `bson:"bookingAndPoliciesPage" json:"bookingAndPoliciesPage"`
	PositionOnMap                 string                        `bson:"positionOnMap" json:"positionOnMap"`
	Contacts                      Contacts                      `bson:"contacts" json:"contacts"`
}
type Contacts struct {
	Phone     string `bson:"phone" json:"phone"`
	Email     string `bson:"email" json:"email"`
	Web       string `bson:"web" json:"web"`
	Facebook  string `bson:"facebook" json:"facebook"`
	Instagram string `bson:"instagram" json:"instagram"`
	Linkedin  string `bson:"linkedin" json:"linkedin"`
}

type OverviewPage struct {
	OverviewText  string          `bson:"overviewText" json:"overviewText"`
	WhyStayWithUs []string        `bson:"whyStayWithUs" json:"whyStayWithUs"`
	QuickFacts    []string        `bson:"quickFacts" json:"quickFacts"`
	Logo          types.FileField `bson:"logo" json:"logo"`
}

type RoomsAndSuitesPage struct {
	StartingText   string         `bson:"startingText" json:"startingText"`
	Advantages     []Advantages   `bson:"advantages" json:"advantages"`
	RoomCategories []RoomCategory `bson:"roomCategories" json:"roomCategories"`
}
type RoomCategory struct {
	RoomType      string   `bson:"roomType" json:"roomType"`
	TotalRoom     int      `bson:"totalRoom" json:"totalRoom"`
	RoomSurface   string   `bson:"roomSurface" json:"roomSurface"`
	BedsCount     int      `bson:"bedsCount" json:"bedsCount"`
	MaxOccupancy  int      `bson:"maxOccupancy" json:"maxOccupancy"`
	ViewType      string   `bson:"viewType" json:"viewType"`
	RoomAmenities []string `bson:"roomAmenities" json:"roomAmenities"`
	Features      []string `bson:"features" json:"features"`
}

type Advantages struct {
	Text string          `bson:"text" json:"text"`
	Icon types.FileField `bson:"icon" json:"icon"`
}

type AmenitiesAndFacilitiesPage struct {
	StartingText                string            `bson:"startingText" json:"startingText"`
	LeisureAndRecreation        []AmenitiesDetail `bson:"leisureAndRecreation" json:"leisureAndRecreation"`
	DiningAndCulinaryExperience []AmenitiesDetail `bson:"diningAndCulinaryExperience" json:"diningAndCulinaryExperience"`
	BusinessAndEvents           []AmenitiesDetail `bson:"businessAndEvents" json:"businessAndEvents"`
	ConvenienceAndServices      []string          `bson:"convenienceAndServices" json:"convenienceAndServices"`
}

type AmenitiesDetail struct {
	Title string          `bson:"title" json:"title"`
	Body  string          `bson:"body" json:"body"`
	Image types.FileField `bson:"image" json:"image"`
}

type LocationNearbyAttractionsPage struct {
	StartingText                   string              `bson:"startingText" json:"startingText"`
	HotelAddress                   HotelAddress        `bson:"hotelAddress" json:"hotelAddress"`
	TopAttractionsNearby           []AttractionsNearby `bson:"topAttractionsNearby" json:"topAttractionsNearby"`
	TransportationAndAccessibility []string            `bson:"transportationAndAccessibility" json:"transportationAndAccessibility"`
}

type HotelAddress struct {
	Longitude   string `bson:"longitude" json:"longitude"`
	Latitude    string `bson:"latitude" json:"latitude"`
	FullAddress string `bson:"fullAddress" json:"fullAddress"`
}

type AttractionsNearby struct {
	Title       string          `bson:"title" json:"title"`
	Body        string          `bson:"body" json:"body"`
	Image       types.FileField `bson:"image" json:"image"`
	Longitude   string          `bson:"longitude" json:"longitude"`
	Latitude    string          `bson:"latitude" json:"latitude"`
	FullAddress string          `bson:"fullAddress" json:"fullAddress"`
}

type ReviewsAndRatingsPage struct {
	StartingText  string             `bson:"startingText" json:"startingText"`
	HotelReviews  []HotelReview      `bson:"hotelReviews" json:"hotelReviews"`
	ActiveProgram primitive.ObjectID `bson:"activeProgram" json:"activeProgram"`
}
type HotelReview struct {
	ReviewerImage types.FileField   `bson:"reviewerImage" json:"reviewerImage"`
	ReviewerName  string            `bson:"reviewerName" json:"reviewerName"`
	ReviewText    string            `bson:"reviewText" json:"reviewText"`
	ImagesGallery []types.FileField `bson:"imagesGallery" json:"imagesGallery"`
}

type BookingAndPoliciesPage struct {
	StartingText                        string   `bson:"startingText" json:"startingText"`
	CheckInAndCheckOut                  []string `bson:"checkInAndCheckOut" json:"checkInAndCheckOut"`
	PaymentPolicies                     []string `bson:"paymentPolicies" json:"paymentPolicies"`
	CancellationPolicy                  []string `bson:"cancellationPolicy" json:"cancellationPolicy"`
	HotelRulesAndPolicies               []string `bson:"hotelRulesAndPolicies" json:"hotelRulesAndPolicies"`
	TransportationAndAdditionalServices []string `bson:"transportationAndAdditionalServices" json:"transportationAndAdditionalServices"`
}

type Hotels struct {
	HotelsDto `bson:",inline"`
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type HotelsRes struct {
	Hotels `bson:",inline"`
	IsFav  bool `bson:"isFav" json:"isFav"`
}

type HotelsPagination struct {
	Hotels     []Hotels          `bson:"hotels" json:"hotels"`
	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}

type CheckInAndCheckOut struct {
	From time.Time `bson:"from" json:"from"`
	To   time.Time `bson:"to" json:"to"`
}
