package filter

import (
	"net/url"
	"strconv"
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
	// Legacy fields for backward compatibility (mapped to correct fields in ToBsonFilter)
	Name     *string    `json:"name,omitempty"`     // Maps to fullName
	Country  *string    `json:"country,omitempty"`  // Maps to nationality
	DateFrom *time.Time `json:"dateFrom,omitempty"` // Maps to createdAt range
	DateTo   *time.Time `json:"dateTo,omitempty"`   // Maps to createdAt range
}

func (f VisitorFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A

	if f.Name != nil && *f.Name != "" {
		ands = append(ands, bson.M{"fullName": bson.M{"$regex": *f.Name, "$options": "i"}})
	}

	if f.Email != "" {
		ands = append(ands, bson.M{"email": bson.M{"$regex": f.Email, "$options": "i"}})
	}

	if f.Phone != "" {
		ands = append(ands, bson.M{"phone": bson.M{"$regex": f.Phone, "$options": "i"}})
	}

	if f.Country != nil && *f.Country != "" {
		ands = append(ands, bson.M{"nationality": bson.M{"$regex": *f.Country, "$options": "i"}})
	}

	dateFilter := bson.M{}
	if f.DateFrom != nil {
		dateFilter["$gte"] = f.DateFrom
	}

	if f.DateTo != nil {
		dateFilter["$lte"] = f.DateTo
	}

	if len(dateFilter) > 0 {
		ands = append(ands, bson.M{"createdAt": dateFilter})
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}

// ParseQueryParams parses URL query parameters into filter values
func (f *VisitorFilter) ParseQueryParams(q url.Values) {
	// Name filter
	if name := q.Get("name"); name != "" {
		f.Name = &name
	}

	// Email filter
	if email := q.Get("email"); email != "" {
		f.Email = email
	}

	// Phone filter
	if phone := q.Get("phone"); phone != "" {
		f.Phone = phone
	}

	// Country filter
	if country := q.Get("country"); country != "" {
		f.Country = &country
	}

	// Note: registrationId and status fields removed as they don't exist in visitor model

	// Date filters
	if from := q.Get("dateFrom"); from != "" {
		date, err := time.Parse("2006-01-02", from)
		if err == nil {
			f.DateFrom = &date
		}
	}

	if to := q.Get("dateTo"); to != "" {
		date, err := time.Parse("2006-01-02", to)
		if err == nil {
			// Set to end of day
			date = date.Add(24*time.Hour - time.Second)
			f.DateTo = &date
		}
	}

	// Set pagination
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	f.Page = page

	size, _ := strconv.Atoi(q.Get("size"))
	if size < 1 {
		size = 10 // Default size
	}
	f.Size = size
}

func (f VisitorFilter) ToBsonFilter() bson.M {
	filter := bson.M{"trash": false}

	// Fix field name: use 'fullName' instead of 'name'
	if f.Name != nil && *f.Name != "" {
		filter["fullName"] = bson.M{"$regex": *f.Name, "$options": "i"}
	}

	// Use correct field name for full name filter
	if f.FullName != "" {
		filter["fullName"] = bson.M{"$regex": f.FullName, "$options": "i"}
	}

	if f.Email != "" {
		filter["email"] = bson.M{"$regex": f.Email, "$options": "i"}
	}

	if f.Phone != "" {
		filter["phone"] = bson.M{"$regex": f.Phone, "$options": "i"}
	}

	// Fix field name: use 'nationality' instead of 'country'
	if f.Country != nil && *f.Country != "" {
		filter["nationality"] = bson.M{"$regex": *f.Country, "$options": "i"}
	}

	// Use correct field name for nationality filter
	if f.Nationality != "" {
		filter["nationality"] = bson.M{"$regex": f.Nationality, "$options": "i"}
	}

	// Handle hotelId filter
	if !f.HotelId.IsZero() {
		filter["hotelId"] = f.HotelId
	}

	// Handle interests filter
	if len(f.Interests) > 0 {
		filter["interests"] = bson.M{"$in": f.Interests}
	}

	// Handle tags filter
	if len(f.Tags) > 0 {
		filter["tags"] = bson.M{"$in": f.Tags}
	}

	// Handle source filter
	if f.Source != "" {
		filter["source"] = bson.M{"$regex": f.Source, "$options": "i"}
	}

	// Handle active status filter
	if f.IsActive != nil {
		filter["isActive"] = *f.IsActive
	}

	// Handle VIP status filter
	if f.IsVIP != nil {
		filter["isVIP"] = *f.IsVIP
	}

	// Handle visit date range
	if f.VisitDateFrom != nil || f.VisitDateTo != nil {
		visitDateFilter := bson.M{}
		if f.VisitDateFrom != nil {
			visitDateFilter["$gte"] = *f.VisitDateFrom
		}
		if f.VisitDateTo != nil {
			visitDateFilter["$lte"] = *f.VisitDateTo
		}
		filter["visitDate"] = visitDateFilter
	}

	// Handle visit count range
	if f.MinVisitCount != nil || f.MaxVisitCount != nil {
		visitCountFilter := bson.M{}
		if f.MinVisitCount != nil {
			visitCountFilter["$gte"] = *f.MinVisitCount
		}
		if f.MaxVisitCount != nil {
			visitCountFilter["$lte"] = *f.MaxVisitCount
		}
		filter["visitCount"] = visitCountFilter
	}

	// Handle total spent range
	if f.MinTotalSpent != nil || f.MaxTotalSpent != nil {
		totalSpentFilter := bson.M{}
		if f.MinTotalSpent != nil {
			totalSpentFilter["$gte"] = *f.MinTotalSpent
		}
		if f.MaxTotalSpent != nil {
			totalSpentFilter["$lte"] = *f.MaxTotalSpent
		}
		filter["totalSpent"] = totalSpentFilter
	}

	// Handle search text (search across multiple fields)
	if f.SearchText != "" {
		searchRegex := bson.M{"$regex": f.SearchText, "$options": "i"}
		filter["$or"] = []bson.M{
			{"fullName": searchRegex},
			{"email": searchRegex},
			{"phone": searchRegex},
			{"nationality": searchRegex},
			{"source": searchRegex},
			{"notes": searchRegex},
		}
	}

	// Handle creation date range
	if f.DateFrom != nil || f.DateTo != nil {
		dateFilter := bson.M{}
		if f.DateFrom != nil {
			dateFilter["$gte"] = *f.DateFrom
		}
		if f.DateTo != nil {
			dateFilter["$lte"] = *f.DateTo
		}
		filter["createdAt"] = dateFilter
	}

	return filter
}
