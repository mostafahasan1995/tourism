package filter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VisitorFilter struct {
	Page          int                `bson:"page" json:"page"`
	Size          int                `bson:"size" json:"size"`
	HotelId       primitive.ObjectID `bson:"hotelId,omitempty" json:"hotelId,omitempty"`
	FullName      string             `bson:"fullName" json:"fullName"`
	Nationality   string             `bson:"nationality" json:"nationality"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone" json:"phone"`
	Interests     []string           `bson:"interests" json:"interests"`
	Tags          []string           `bson:"tags" json:"tags"`
	Source        string             `bson:"source" json:"source"`
	IsActive      *bool              `bson:"isActive,omitempty" json:"isActive,omitempty"`
	IsVIP         *bool              `bson:"isVIP,omitempty" json:"isVIP,omitempty"`
	VisitDateFrom *time.Time         `bson:"visitDateFrom,omitempty" json:"visitDateFrom,omitempty"`
	VisitDateTo   *time.Time         `bson:"visitDateTo,omitempty" json:"visitDateTo,omitempty"`
	MinVisitCount *int               `bson:"minVisitCount,omitempty" json:"minVisitCount,omitempty"`
	MaxVisitCount *int               `bson:"maxVisitCount,omitempty" json:"maxVisitCount,omitempty"`
	MinTotalSpent *float64           `bson:"minTotalSpent,omitempty" json:"minTotalSpent,omitempty"`
	MaxTotalSpent *float64           `bson:"maxTotalSpent,omitempty" json:"maxTotalSpent,omitempty"`
	SearchText    string             `bson:"searchText" json:"searchText"`
}

func (f *VisitorFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add hotel ID filter if provided
	if !f.HotelId.IsZero() {
		filterConditions = append(filterConditions, bson.M{"hotelId": f.HotelId})
	}

	// Add full name search if provided
	if f.FullName != "" {
		filterConditions = append(filterConditions, bson.M{
			"fullName": bson.M{"$regex": f.FullName, "$options": "i"},
		})
	}

	// Add nationality filter if provided
	if f.Nationality != "" {
		filterConditions = append(filterConditions, bson.M{
			"nationality": bson.M{"$regex": f.Nationality, "$options": "i"},
		})
	}

	// Add email search if provided
	if f.Email != "" {
		filterConditions = append(filterConditions, bson.M{
			"email": bson.M{"$regex": f.Email, "$options": "i"},
		})
	}

	// Add phone search if provided
	if f.Phone != "" {
		filterConditions = append(filterConditions, bson.M{
			"phone": bson.M{"$regex": f.Phone, "$options": "i"},
		})
	}

	// Add interests filter if provided
	if len(f.Interests) > 0 {
		filterConditions = append(filterConditions, bson.M{
			"interests": bson.M{"$in": f.Interests},
		})
	}

	// Add tags filter if provided
	if len(f.Tags) > 0 {
		filterConditions = append(filterConditions, bson.M{
			"tags": bson.M{"$in": f.Tags},
		})
	}

	// Add source filter if provided
	if f.Source != "" {
		filterConditions = append(filterConditions, bson.M{"source": f.Source})
	}

	// Add active status filter if provided
	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	// Add VIP status filter if provided
	if f.IsVIP != nil {
		filterConditions = append(filterConditions, bson.M{"isVIP": *f.IsVIP})
	}

	// Add visit date range filter if provided
	if f.VisitDateFrom != nil && f.VisitDateTo != nil {
		filterConditions = append(filterConditions, bson.M{
			"visitDate": bson.M{
				"$gte": *f.VisitDateFrom,
				"$lte": *f.VisitDateTo,
			},
		})
	} else if f.VisitDateFrom != nil {
		filterConditions = append(filterConditions, bson.M{
			"visitDate": bson.M{"$gte": *f.VisitDateFrom},
		})
	} else if f.VisitDateTo != nil {
		filterConditions = append(filterConditions, bson.M{
			"visitDate": bson.M{"$lte": *f.VisitDateTo},
		})
	}

	// Add visit count range filter if provided
	if f.MinVisitCount != nil && f.MaxVisitCount != nil {
		filterConditions = append(filterConditions, bson.M{
			"visitCount": bson.M{
				"$gte": *f.MinVisitCount,
				"$lte": *f.MaxVisitCount,
			},
		})
	} else if f.MinVisitCount != nil {
		filterConditions = append(filterConditions, bson.M{
			"visitCount": bson.M{"$gte": *f.MinVisitCount},
		})
	} else if f.MaxVisitCount != nil {
		filterConditions = append(filterConditions, bson.M{
			"visitCount": bson.M{"$lte": *f.MaxVisitCount},
		})
	}

	// Add total spent range filter if provided
	if f.MinTotalSpent != nil && f.MaxTotalSpent != nil {
		filterConditions = append(filterConditions, bson.M{
			"totalSpent": bson.M{
				"$gte": *f.MinTotalSpent,
				"$lte": *f.MaxTotalSpent,
			},
		})
	} else if f.MinTotalSpent != nil {
		filterConditions = append(filterConditions, bson.M{
			"totalSpent": bson.M{"$gte": *f.MinTotalSpent},
		})
	} else if f.MaxTotalSpent != nil {
		filterConditions = append(filterConditions, bson.M{
			"totalSpent": bson.M{"$lte": *f.MaxTotalSpent},
		})
	}

	// Add text search if provided
	if f.SearchText != "" {
		filterConditions = append(filterConditions, bson.M{
			"$or": []bson.M{
				{"fullName": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"email": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"phone": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"nationality": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"interests": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"tags": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"notes": bson.M{"$regex": f.SearchText, "$options": "i"}},
			},
		})
	}

	return bson.M{"$and": filterConditions}
}
