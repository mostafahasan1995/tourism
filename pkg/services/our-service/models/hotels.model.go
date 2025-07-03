package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	//"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelsDto struct {
	Name                          transl.Localizable[string]    `bson:"name" json:"name"`
	HotelType                     string                        `bson:"hotelType" json:"hotelType"`
	Location                      string                        `bson:"location" json:"location"`
	CheckInAndCheckOut            CheckInAndCheckOut            `bson:"checkInAndCheckOut" json:"checkInAndCheckOut"`
	Price                         int                           `bson:"price" json:"price"`
	IsDisplayInPerfectStay        bool                          `bson:"isDisplayInPerfectStay" json:"isDisplayInPerfectStay"`
	Ratings                       float64                       `bson:"ratings" json:"ratings"`
	Image                         types.FileField               `bson:"image" json:"image"`
	RoomAmenities                 []transl.Localizable[string]  `bson:"roomAmenities" json:"roomAmenities"`
	DistanceFromCityCenter        int                           `bson:"distanceFromCityCenter" json:"distanceFromCityCenter"`
	NearbyAttractions             []transl.Localizable[string]  `bson:"nearbyAttractions" json:"nearbyAttractions"`
	ImagesGallery                 []types.FileField             `bson:"imagesGallery" json:"imagesGallery"`
	OverviewPage                  OverviewPage                  `bson:"overviewPage" json:"overviewPage"`
	RoomsAndSuitesPage            RoomsAndSuitesPage            `bson:"roomsAndSuitesPage" json:"roomsAndSuitesPage"`
	AmenitiesAndFacilitiesPage    AmenitiesAndFacilitiesPage    `bson:"amenitiesAndFacilitiesPage" json:"amenitiesAndFacilitiesPage"`
	AmenitiesAndFacilitiesPageV2  AmenitiesAndFacilitiesPageV2  `bson:"amenitiesAndFacilitiesPageV2" json:"amenitiesAndFacilitiesPageV2"`
	LocationNearbyAttractionsPage LocationNearbyAttractionsPage `bson:"locationNearbyAttractionsPage" json:"locationNearbyAttractionsPage"`
	ReviewsAndRatingsPage         ReviewsAndRatingsPage         `bson:"reviewsAndRatingsPage" json:"reviewsAndRatingsPage"`
	BookingAndPoliciesPage        BookingAndPoliciesPage        `bson:"bookingAndPoliciesPage" json:"bookingAndPoliciesPage"`
	PositionOnMap                 string                        `bson:"positionOnMap" json:"positionOnMap"`
	Contacts                      Contacts                      `bson:"contacts" json:"contacts"`
	CloseReservations             []CloseReservations           `bson:"closeReservations" json:"closeReservations"`
	RatingObjects                 []RatingObject                `bson:"ratingObjects" json:"ratingObjects"`
	OfferAndDiscount              []OfferAndDiscount            `bson:"offerAndDiscount" json:"offerAndDiscount"`
	PoliciesPage                  PoliciesPage                  `bson:"policiesPage" json:"policiesPage"`
	Owner                         primitive.ObjectID            `bson:"owner" json:"owner"`
	Country                       string                        `bson:"country" json:"country"`
}

// CalculateAverageRating calculates the average rating from RatingObjects
func (h *HotelsDto) CalculateAverageRating() {
	if len(h.RatingObjects) == 0 {
		return
	}
	var total float64
	for _, rating := range h.RatingObjects {
		total += rating.Value
	}
	h.Ratings = total / float64(len(h.RatingObjects))
}

type RatingObject struct {
	Username string          `bson:"username" json:"username"`
	UserId   string          `bson:"userId" json:"userId"`
	UserImg  types.FileField `bson:"userImg" json:"userImg"`
	Value    float64         `bson:"value" json:"value"`
	Text     string          `bson:"text" json:"text"`
	Status   string          `bson:"status" json:"status"`
	Date     *time.Time      `bson:"date" json:"date"`
	Replies  []Reply         `bson:"replies" json:"replies"`
}
type Reply struct {
	Text string     `bson:"text" json:"text"`
	Date *time.Time `bson:"date" json:"date"`
}

type CloseReservations struct {
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
}

type OfferAndDiscount struct {
	Name        string    `bson:"name" json:"name"`
	Code        string    `bson:"code" json:"code"`
	StartDate   time.Time `bson:"startDate" json:"startDate"`
	EndDate     time.Time `bson:"endDate" json:"endDate"`
	Description string    `bson:"description" json:"description"`
}

type Phone struct {
	Pre     string `bson:"pre" json:"pre"`
	Content string `bson:"content" json:"content"`
}
type Contacts struct {
	Phone  Phone    `bson:"phone" json:"phone"`
	Email  string   `bson:"email" json:"email"`
	Web    string   `bson:"web" json:"web"`
	Social []Social `bson:"social" json:"social"`
}
type Social struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}

type OverviewPage struct {
	OverviewText  transl.Localizable[string]   `bson:"overviewText" json:"overviewText"`
	WhyStayWithUs []transl.Localizable[string] `bson:"whyStayWithUs" json:"whyStayWithUs"`
	QuickFacts    []transl.Localizable[string] `bson:"quickFacts" json:"quickFacts"`
	Logo          types.FileField              `bson:"logo" json:"logo"`
}

type RoomsAndSuitesPage struct {
	StartingText   transl.Localizable[string]   `bson:"startingText" json:"startingText"`
	Advantages     []transl.Localizable[string] `bson:"advantages" json:"advantages"`
	RoomCategories []RoomCategory               `bson:"roomCategories" json:"roomCategories"`
}
type RoomCategory struct {
	RoomType          string                       `bson:"roomType" json:"roomType"`
	TotalRoom         int                          `bson:"totalRoom" json:"totalRoom"`
	RoomSurface       string                       `bson:"roomSurface" json:"roomSurface"`
	BedsCount         int                          `bson:"bedsCount" json:"bedsCount"`
	MaxOccupancy      int                          `bson:"maxOccupancy" json:"maxOccupancy"`
	ViewType          string                       `bson:"viewType" json:"viewType"`
	ActivePricingType string                       `bson:"activePricingType" json:"activePricingType"`
	RoomAmenities     []string                     `bson:"roomAmenities" json:"roomAmenities"`
	Features          []transl.Localizable[string] `bson:"features" json:"features"`
	Images            []types.FileField            `bson:"images" json:"images"`
	Pricing           Pricing                      `bson:"pricing" json:"pricing"`
	SeasonalPricing   []SeasonPricing              `bson:"seasonalPricing" json:"seasonalPricing"`
}

type Pricing struct {
	NightlyRateBase       int    `bson:"nightlyRateBase" json:"nightlyRateBase"`
	NightlyRateBaseType   string `bson:"nightlyRateBaseType" json:"nightlyRateBaseType"`
	ExtraPersonCharge     int    `bson:"extraPersonCharge" json:"extraPersonCharge"`
	ExtraPersonChargeType string `bson:"extraPersonChargeType" json:"extraPersonChargeType"`
	IsIncludeBreakFast    bool   `bson:"isIncludeBreakFast" json:"isIncludeBreakFast"`

	CurrencyType string `bson:"currencyType" json:"currencyType"`
}
type SeasonPricing struct {
	NightlyRateBase       int    `bson:"nightlyRateBase" json:"nightlyRateBase"`
	NightlyRateBaseType   string `bson:"nightlyRateBaseType" json:"nightlyRateBaseType"`
	ExtraPersonCharge     int    `bson:"extraPersonCharge" json:"extraPersonCharge"`
	ExtraPersonChargeType string `bson:"extraPersonChargeType" json:"extraPersonChargeType"`
	IsIncludeBreakFast    bool   `bson:"isIncludeBreakFast" json:"isIncludeBreakFast"`
	SeasonName            string `bson:"seasonName" json:"seasonName"`

	CurrencyType string `bson:"currencyType" json:"currencyType"`
}
type Advantages struct {
	Text string          `bson:"text" json:"text"`
	Icon types.FileField `bson:"icon" json:"icon"`
}

type AmenitiesAndFacilitiesPage struct {
	StartingText                transl.Localizable[string]   `bson:"startingText" json:"startingText"`
	LeisureAndRecreation        []AmenitiesDetail            `bson:"leisureAndRecreation" json:"leisureAndRecreation"`
	DiningAndCulinaryExperience []AmenitiesDetail            `bson:"diningAndCulinaryExperience" json:"diningAndCulinaryExperience"`
	BusinessAndEvents           []AmenitiesDetail            `bson:"businessAndEvents" json:"businessAndEvents"`
	ConvenienceAndServices      []transl.Localizable[string] `bson:"convenienceAndServices" json:"convenienceAndServices"`
}

type AmenitiesAndFacilitiesPageV2 struct {
	RestaurantsCafes       []string `bson:"restaurantsCafes" json:"restaurantsCafes"`
	PoolsBeaches           []string `bson:"poolsBeaches" json:"poolsBeaches"`
	SpaGym                 []string `bson:"spaGym" json:"spaGym"`
	HotelServices          []string `bson:"hotelServices" json:"hotelServices"`
	BusinessFacilities     []string `bson:"businessFacilities" json:"businessFacilities"`
	KidsFacilities         []string `bson:"kidsFacilities" json:"kidsFacilities"`
	RecreationalActivities []string `bson:"recreationalActivities" json:"recreationalActivities"`
}

type AmenitiesDetail struct {
	Title transl.Localizable[string] `bson:"title" json:"title"`
	Body  transl.Localizable[string] `bson:"body" json:"body"`
	Image types.FileField            `bson:"image" json:"image"`
}

type LocationNearbyAttractionsPage struct {
	StartingText                   transl.Localizable[string]       `bson:"startingText" json:"startingText"`
	HotelAddress                   HotelAddress                     `bson:"hotelAddress" json:"hotelAddress"`
	TopAttractionsNearby           []AttractionsNearby              `bson:"topAttractionsNearby" json:"topAttractionsNearby"`
	TransportationAndAccessibility []TransportationAndAccessibility `bson:"transportationAndAccessibility" json:"transportationAndAccessibility"`
}

type TransportationAndAccessibility struct {
	Title        string `bson:"title" json:"title"`
	Description  string `bson:"description" json:"description"`
	Availability bool   `bson:"availability" json:"availability"`
}

type HotelAddress struct {
	Longitude   string                     `bson:"longitude" json:"longitude"`
	Latitude    string                     `bson:"latitude" json:"latitude"`
	FullAddress transl.Localizable[string] `bson:"fullAddress" json:"fullAddress"`
}

type AttractionsNearby struct {
	Title       transl.Localizable[string] `bson:"title" json:"title"`
	Body        transl.Localizable[string] `bson:"body" json:"body"`
	Image       types.FileField            `bson:"image" json:"image"`
	Longitude   string                     `bson:"longitude" json:"longitude"`
	Latitude    string                     `bson:"latitude" json:"latitude"`
	FullAddress transl.Localizable[string] `bson:"fullAddress" json:"fullAddress"`
}

type ReviewsAndRatingsPage struct {
	StartingText  transl.Localizable[string] `bson:"startingText" json:"startingText"`
	HotelReviews  []HotelReviewDisplay       `bson:"hotelReviews" json:"hotelReviews"`
	ActiveProgram primitive.ObjectID         `bson:"activeProgram" json:"activeProgram"`
}

type HotelReviewDisplay struct {
	ReviewerImage types.FileField            `bson:"reviewerImage" json:"reviewerImage"`
	ReviewerName  transl.Localizable[string] `bson:"reviewerName" json:"reviewerName"`
	ReviewText    transl.Localizable[string] `bson:"reviewText" json:"reviewText"`
	ImagesGallery []types.FileField          `bson:"imagesGallery" json:"imagesGallery"`
}

type BookingAndPoliciesPage struct {
	StartingText                        transl.Localizable[string] `bson:"startingText" json:"startingText"`
	CheckInAndCheckOut                  []string                   `bson:"checkInAndCheckOut" json:"checkInAndCheckOut"`
	PaymentPolicies                     []string                   `bson:"paymentPolicies" json:"paymentPolicies"`
	CancellationPolicy                  []string                   `bson:"cancellationPolicy" json:"cancellationPolicy"`
	HotelRulesAndPolicies               []string                   `bson:"hotelRulesAndPolicies" json:"hotelRulesAndPolicies"`
	TransportationAndAdditionalServices []string                   `bson:"transportationAndAdditionalServices" json:"transportationAndAdditionalServices"`
}
type PoliciesPage struct {
	CheckInCheckOut                  *map[string]any `bson:"checkInCheckOut,omitempty" json:"checkInCheckOut,omitempty"`
	PaymentPolicies                  *map[string]any `bson:"paymentPolicies,omitempty" json:"paymentPolicies,omitempty"`
	CancellationPolicy               *map[string]any `bson:"cancellationPolicy,omitempty" json:"cancellationPolicy,omitempty"`
	HotelRulesPolicies               *map[string]any `bson:"hotelRulesPolicies,omitempty" json:"hotelRulesPolicies,omitempty"`
	TransportationAdditionalServices *map[string]any `bson:"transportationAdditionalServices,omitempty" json:"transportationAdditionalServices,omitempty"`
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
