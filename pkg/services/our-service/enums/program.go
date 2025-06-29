package enums

//travl request service type
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
)

type ProgramServiceType string

const (
	ProgramServiceTypeTourismProgram ProgramServiceType = "tourism-program"
	ProgramServiceTypeCustomProgram  ProgramServiceType = "custom-program"
	ProgramServiceTypeFlightTicket   ProgramServiceType = "flight-ticket"
	ProgramServiceTypeVipCar         ProgramServiceType = "vip-car"
	ProgramServiceTypeHotelBooking   ProgramServiceType = "hotel-booking"
)

//travel type
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
)

//program type

type ProgramType string

const (
	ProgramTypeGeneral ProgramType = "general"
	ProgramTypeCustom  ProgramType = "custom"
)
