package enums

// travl request service type
type ServiceType string

const (
	// Basic Service Types
	ServiceTypeDelegation     ServiceType = "delegation"
	ServiceTypeCustomPlan     ServiceType = "custom-plan"
	ServiceTypeBusinessMan    ServiceType = "business-man"
	ServiceTypeVipCar         ServiceType = "vip-car"
	ServiceTypeFlightRequest  ServiceType = "flight-request"
	ServiceTypePartnerRequest ServiceType = "partner-request"
	ServiceTypeHotelBooking   ServiceType = "hotel-booking"
)

type GroupSize string

const (
	// Group Size Options (عدد الأفراد)
	GroupSizeSolo   GroupSize = "solo"   // Solo Traveler (فردي)
	GroupSizeCouple GroupSize = "couple" // Couple (زوجان)
	GroupSizeFamily GroupSize = "family" // Family (عائلة)
	GroupSizeSmall  GroupSize = "small"  // Small Group (مجموعة صغيرة، عادة 4–8 أشخاص)
	GroupSizeLarge  GroupSize = "large"  // Large Group (مجموعة كبيرة، عادة أكثر من 8 أشخاص)
	GroupSizeCustom GroupSize = "custom" // Custom Group (مجموعة خاصة)
	GroupSizeFixed  GroupSize = "fixed"  // Fixed Group (مجموعة ثابتة)
	GroupSizeMedium GroupSize = "medium" // Medium Group (مجموعة متوسطة)
)

type ProgramServiceType string

const (
	ProgramServiceTypeTourismProgram    ProgramServiceType = "tourism-program"
	ProgramServiceTypeCustomProgram     ProgramServiceType = "custom-program"
	ProgramServiceTypeFlightTicket      ProgramServiceType = "flight-ticket"
	ProgramServiceTypeVipCar            ProgramServiceType = "vip-car"
	ProgramServiceTypeHotelBooking      ProgramServiceType = "hotel-booking"
	ProgramServiceTypeFamilyTravel      ProgramServiceType = "family-travel"
	ProgramServiceTypeLuxuryTravel      ProgramServiceType = "luxury-travel"
	ProgramServiceTypeReligiousTravel   ProgramServiceType = "religious-travel"
	ProgramServiceTypeHoneymoon         ProgramServiceType = "honeymoon"
	ProgramServiceTypeBusinessManTravel ProgramServiceType = "business-man-travel"
	ProgramServiceTypeDelegation        ProgramServiceType = "delegation"
)

// travel type
type TravelType string

const (

	// Travel Type Options
	TravelTypeRelaxationTrip      TravelType = "relaxation-trip"          // Relaxation Trip (رحلة استجمام)
	TravelTypeAdventureTrip       TravelType = "adventure-trip"           // Adventure (مغامرة)
	TravelTypeFamilyTrip          TravelType = "family-trip"              // Family Trip (رحلة عائلية)
	TravelTypeRomanticTrip        TravelType = "romantic-trip"            // Romantic Trip (Honeymoon) (رحلة رومانسية – شهر عسل)
	TravelTypeCulturalTrip        TravelType = "cultural-trip"            // Cultural Trip (رحلة ثقافية)
	TravelTypeBusinessTrip        TravelType = "business-trip"            // Business Trip (رحلة عمل)
	TravelTypeShoppingTrip        TravelType = "shopping-trip"            // Shopping Trip (رحلة تسوق)
	TravelTypeWellnessMedicalTrip TravelType = "wellness-medical-tourism" // Wellness or Medical Tourism (رحلة صحية أو استشفائية)

	// Additional Travel Types
	TravelTypeLeisureTravel                 TravelType = "leisure-travel"                   // Leisure Travel
	TravelTypeAdventureTravel               TravelType = "adventure-travel"                 // Adventure Travel
	TravelTypeLuxuryTravel                  TravelType = "luxury-travel"                    // Luxury Travel
	TravelTypeCulturalTravel                TravelType = "cultural-travel"                  // Cultural Travel
	TravelTypeNatureWildlifeTravel          TravelType = "nature-wildlife-travel"           // Nature & Wildlife Travel
	TravelTypeReligiousTravel               TravelType = "religious-travel"                 // Religious Travel
	TravelTypeRomanticTravel                TravelType = "romantic-travel"                  // Romantic Travel
	TravelTypeFamilyTravel                  TravelType = "family-travel"                    // Family Travel
	TravelTypeEcoTravel                     TravelType = "eco-travel"                       // Eco Travel
	TravelTypeCruiseTravel                  TravelType = "cruise-travel"                    // Cruise Travel
	TravelTypeMultiCountryTravel            TravelType = "multi-country-travel"             // Multi-Country Travel
	TravelTypeCityBreakTravel               TravelType = "city-break-travel"                // City Break Travel
	TravelTypeWellnessTravel                TravelType = "wellness-travel"                  // Wellness Travel
	TravelTypeEventBasedTravel              TravelType = "event-based-travel"               // Event-Based Travel
	TravelTypeVipCelebrityTravel            TravelType = "vip-celebrity-travel"             // VIP / Celebrity Travel
	TravelTypeLuxuryEscape                  TravelType = "luxury-escape"                    // Luxury Escape
	TravelTypeHoneymoon                     TravelType = "honeymoon"                        // Honeymoon
	TravelTypeFamilyLuxuryHoliday           TravelType = "family-luxury-holiday"            // Family Luxury Holiday
	TravelTypeWellnessSpaRetreat            TravelType = "wellness-spa-retreat"             // Wellness Spa Retreat
	TravelTypeCulinaryFineDiningTour        TravelType = "culinary-fine-dining-tour"        // Culinary Fine Dining Tour
	TravelTypeLuxuryCruiseExperience        TravelType = "luxury-cruise-experience"         // Luxury Cruise Experience
	TravelTypeWinterSkiRetreat              TravelType = "winter-ski-retreat"               // Winter Ski Retreat
	TravelTypeCustomVipTour                 TravelType = "custom-vip-tour"                  // Custom VIP Tour
	TravelTypeExclusiveSafariNature         TravelType = "exclusive-safari-nature"          // Exclusive Safari Nature
	TravelTypeBeachIslandGetaway            TravelType = "beach-island-getaway"             // Beach Island Getaway
	TravelTypeCulturalHeritageJourney       TravelType = "cultural-heritage-journey"        // Cultural Heritage Journey
	TravelTypeReligiousSpiritualJourney     TravelType = "religious-spiritual-journey"      // Religious Spiritual Journey
	TravelTypeIconicLandmarksCityHighlights TravelType = "iconic-landmarks-city-highlights" // Iconic Landmarks City Highlights
	TravelTypePrivateGuidedTour             TravelType = "private-guided-tour"              // Private Guided Tour
	TravelTypeCouplesPrivateTour            TravelType = "couples-private-tour"             // Couples Private Tour
	TravelTypeFestivalSpecialEvents         TravelType = "festival-special-events"          // Festival Special Events
	TravelTypePhotographyScenicTour         TravelType = "photography-scenic-tour"          // Photography Scenic Tour
	TravelTypeMiceBusinessTravel            TravelType = "mice-business-travel"             // Mice Business Travel

)

//program type

type ProgramType string

const (
	ProgramTypeGeneral ProgramType = "general"
	ProgramTypeCustom  ProgramType = "custom"
)
