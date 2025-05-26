package models

import (
	"larsa-tourism-microservices/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Travel Type values:
// Relaxation Trip (رحلة استجمام)
// Adventure (مغامرة)
// Family Trip (رحلة عائلية)
// Romantic Trip (Honeymoon) (رحلة رومانسية – شهر عسل)
// Cultural Trip (رحلة ثقافية)
// Business Trip (رحلة عمل)
// Shopping Trip (رحلة تسوق)
// Wellness or Medical Tourism (رحلة صحية أو استشفائية)

// ✅ عدد الأفراد (Group Size):
// Solo Traveler (فردي)s
// Couple (زوجان)
// Family (عائلة)
// Small Group (مجموعة صغيرة، عادة 4–8 أشخاص)
// Large Group (مجموعة كبيرة، عادة أكثر من 8 أشخاص)

// general program - to show in the website so the customer can book it then we will create a custom program for them
type GeneralProgram struct {
	//BasicInfo      BasicInfo            `bson:"basicInfo" json:"basicInfo"`
	Distinations   []primitive.ObjectID `bson:"distinations" json:"distinations"`
	Includes       Includes             `bson:"includes" json:"includes"`
	Activities     []primitive.ObjectID `bson:"activities" json:"activities"`
	DailyItinerary []DailyItinerary     `bson:"dailyItinerary" json:"dailyItinerary"`
	Pricing        GPPricing            `bson:"pricing" json:"pricing"`
}

type Includes struct {
	Accommodation  []string `bson:"accommodation" json:"accommodation"`
	Transportation []string `bson:"transportation" json:"transportation"`
	Meals          []string `bson:"meals" json:"meals"`
}

type DailyItinerary struct {
	Title   string               `bson:"title" json:"title"`
	Actions []primitive.ObjectID `bson:"actions" json:"actions"`
	Images  []types.FileField    `bson:"images" json:"images"`
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

// type GeneralProgram struct {
// 	Id                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
// 	GeneralProgramDto `bson:",inline"`
// 	Trash             bool               `bson:"trash" json:"trash"`
// 	CreatedAt         time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
// 	CreatedBy         primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
// 	UpdatedAt         time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
// 	UpdatedBy         primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
// }

// type GeneralProgramsPagination struct {
// 	Programs   []GeneralProgram `bson:"programs" json:"programs"`
// 	Pagination types.Pagination `bson:"pagination" json:"pagination"`
// }
