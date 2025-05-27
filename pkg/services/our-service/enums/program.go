package enums

type ServiceType string

const (
	// Basic Service Types
	ServiceTypeDelegation     ServiceType = "delegation"
	ServiceTypeCustomPlan     ServiceType = "custom-plan"
	ServiceTypeBusinessMan    ServiceType = "business-man"
	ServiceTypeVipCar         ServiceType = "vip-car"
	ServiceTypeFlightRequest  ServiceType = "flight-request"
	ServiceTypePartnerRequest ServiceType = "partner-request"

	// Travel Types
	ServiceTypeRelaxation ServiceType = "relaxation" // رحلة استجمام
	ServiceTypeAdventure  ServiceType = "adventure"  // مغامرة
	ServiceTypeFamily     ServiceType = "family"     // رحلة عائلية
	ServiceTypeRomantic   ServiceType = "romantic"   // رحلة رومانسية – شهر عسل
	ServiceTypeCultural   ServiceType = "cultural"   // رحلة ثقافية
	ServiceTypeBusiness   ServiceType = "business"   // رحلة عمل
	ServiceTypeShopping   ServiceType = "shopping"   // رحلة تسوق
	ServiceTypeWellness   ServiceType = "wellness"   // رحلة صحية أو استشفائية
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
